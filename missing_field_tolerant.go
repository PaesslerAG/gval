package gval

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// MissingFieldTolerantLogic creates a language extension that treats missing fields as false
// in logical expressions (&&, ||) instead of throwing an error.
func MissingFieldTolerantLogic() Language {
	return NewLanguage(
		// Override the && operator to handle missing field errors
		InfixShortCircuit("&&", func(a interface{}) (interface{}, bool) { 
			return false, a == false 
		}),
		InfixEvalOperator("&&", func(a, b Evaluable) (Evaluable, error) {
			return func(c context.Context, v interface{}) (interface{}, error) {
				// Evaluate left operand
				aVal, err := a(c, v)
				if err != nil {
					// Check if this is a missing field error
					if isMissingFieldError(err) {
						// Treat missing field as false, so AND operation is false
						return false, nil
					}
					return nil, err
				}
				
				// Short circuit if left operand is false
				aBool, ok := convertToBool(aVal)
				if !ok {
					return nil, err
				}
				if !aBool {
					return false, nil
				}
				
				// Evaluate right operand
				bVal, err := b(c, v)
				if err != nil {
					// Check if this is a missing field error
					if isMissingFieldError(err) {
						// Treat missing field as false, so AND operation is false
						return false, nil
					}
					return nil, err
				}
				
				bBool, ok := convertToBool(bVal)
				if !ok {
					return false, nil
				}
				
				return aBool && bBool, nil
			}, nil
		}),
		
		// Override the || operator to handle missing field errors
		InfixShortCircuit("||", func(a interface{}) (interface{}, bool) { 
			return true, a == true 
		}),
		InfixEvalOperator("||", func(a, b Evaluable) (Evaluable, error) {
			return func(c context.Context, v interface{}) (interface{}, error) {
				// Evaluate left operand
				aVal, err := a(c, v)
				if err != nil {
					// Check if this is a missing field error
					if isMissingFieldError(err) {
						// Treat missing field as false, continue to right operand
						aVal = false
					} else {
						return nil, err
					}
				}
				
				// Short circuit if left operand is true
				aBool, ok := convertToBool(aVal)
				if !ok {
					aBool = false
				}
				if aBool {
					return true, nil
				}
				
				// Evaluate right operand
				bVal, err := b(c, v)
				if err != nil {
					// Check if this is a missing field error
					if isMissingFieldError(err) {
						// Treat missing field as false
						return aBool || false, nil
					}
					return nil, err
				}
				
				bBool, ok := convertToBool(bVal)
				if !ok {
					return false, nil
				}
				
				return aBool || bBool, nil
			}, nil
		}),
	)
}

// MissingFieldAsNil creates a language extension that treats missing fields as nil
// instead of throwing an error. This allows expressions to continue evaluation.
func MissingFieldAsNil() Language {
	return VariableSelector(func(path Evaluables) Evaluable {
		return func(c context.Context, v interface{}) (interface{}, error) {
			keys, err := path.EvalStrings(c, v)
			if err != nil {
				return nil, err
			}
			for _, k := range keys {
				switch o := v.(type) {
				case Selector:
					v, err = o.SelectGVal(c, k)
					if err != nil {
						return nil, fmt.Errorf("failed to select '%s' on %T: %w", k, o, err)
					}
					continue
				case map[interface{}]interface{}:
					if val, exists := o[k]; exists {
						v = val
					} else {
						return nil, nil // Return nil instead of error for missing field
					}
					continue
				case map[string]interface{}:
					if val, exists := o[k]; exists {
						v = val
					} else {
						return nil, nil // Return nil instead of error for missing field
					}
					continue
				case []interface{}:
					if i, err := strconv.Atoi(k); err == nil && i >= 0 && len(o) > i {
						v = o[i]
						continue
					}
					return nil, nil // Return nil instead of error for missing array index
				default:
					var ok bool
					v, ok = reflectSelect(k, o)
					if !ok {
						return nil, nil // Return nil instead of error for missing field
					}
				}
			}
			return v, nil
		}
	})
}

// NilSafeComparison creates operators that handle nil values gracefully
func NilSafeComparison() Language {
	return NewLanguage(
		// Override comparison operators to handle nil gracefully
		InfixOperator(">", func(a, b interface{}) (interface{}, error) {
			if a == nil || b == nil {
				return false, nil
			}
			// Try numeric comparison first
			if aFloat, aOk := convertToFloat(a); aOk {
				if bFloat, bOk := convertToFloat(b); bOk {
					return aFloat > bFloat, nil
				}
			}
			// Fall back to string comparison
			return fmt.Sprintf("%v", a) > fmt.Sprintf("%v", b), nil
		}),
		InfixOperator(">=", func(a, b interface{}) (interface{}, error) {
			if a == nil || b == nil {
				return false, nil
			}
			if aFloat, aOk := convertToFloat(a); aOk {
				if bFloat, bOk := convertToFloat(b); bOk {
					return aFloat >= bFloat, nil
				}
			}
			return fmt.Sprintf("%v", a) >= fmt.Sprintf("%v", b), nil
		}),
		InfixOperator("<", func(a, b interface{}) (interface{}, error) {
			if a == nil || b == nil {
				return false, nil
			}
			if aFloat, aOk := convertToFloat(a); aOk {
				if bFloat, bOk := convertToFloat(b); bOk {
					return aFloat < bFloat, nil
				}
			}
			return fmt.Sprintf("%v", a) < fmt.Sprintf("%v", b), nil
		}),
		InfixOperator("<=", func(a, b interface{}) (interface{}, error) {
			if a == nil || b == nil {
				return false, nil
			}
			if aFloat, aOk := convertToFloat(a); aOk {
				if bFloat, bOk := convertToFloat(b); bOk {
					return aFloat <= bFloat, nil
				}
			}
			return fmt.Sprintf("%v", a) <= fmt.Sprintf("%v", b), nil
		}),
		InfixOperator("==", func(a, b interface{}) (interface{}, error) {
			if a == nil && b == nil {
				return true, nil
			}
			if a == nil || b == nil {
				return false, nil
			}
			return reflect.DeepEqual(a, b), nil
		}),
		InfixOperator("!=", func(a, b interface{}) (interface{}, error) {
			if a == nil && b == nil {
				return false, nil
			}
			if a == nil || b == nil {
				return true, nil
			}
			return !reflect.DeepEqual(a, b), nil
		}),
	)
}

// isMissingFieldError checks if an error is due to a missing field
func isMissingFieldError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "unknown parameter")
}
