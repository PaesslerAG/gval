# Missing Field Handling in Gval

This solution addresses the issue where gval expressions fail when trying to access missing nested fields, instead of treating them as false in logical expressions.

## Problem

By default, gval throws an "unknown parameter" error when accessing missing fields:

```go
params := map[string]interface{}{
    "foo": 10,
    "bar": "baz",
    "foo1": map[string]interface{}{"xyz1": 100},
    // Note: foo1.xyz is missing
}

// These expressions will fail with "unknown parameter foo1.xyz" error:
result, err := gval.Evaluate("foo1.xyz > 5", params)                    // ❌ Error
result, err := gval.Evaluate("foo1.xyz > 5 && bar == \"baz\"", params) // ❌ Error  
result, err := gval.Evaluate("foo1.xyz > 5 || bar == \"baz\"", params) // ❌ Error
```

## Solution

The solution provides two new language extensions that work together:

### 1. `MissingFieldAsNil()` - Variable Selector Extension

This extension modifies the variable resolution to return `nil` instead of throwing an error when a field is missing.

### 2. `NilSafeComparison()` - Comparison Operators Extension  

This extension modifies comparison operators (`>`, `>=`, `<`, `<=`, `==`, `!=`) to handle `nil` values gracefully:
- `nil` compared to any value returns `false` (except `nil == nil` which returns `true`)
- `nil != any_value` returns `true` (except `nil != nil` which returns `false`)

## Usage

```go
import "github.com/Nandagopi/gval"

// Create a language that handles missing fields
lang := gval.Full(gval.MissingFieldAsNil(), gval.NilSafeComparison())

params := map[string]interface{}{
    "foo": 10,
    "bar": "baz", 
    "foo1": map[string]interface{}{"xyz1": 100},
}

// Now these expressions work as expected:
result, _ := gval.Evaluate("foo1.xyz > 5", params, lang)                    // ✅ false
result, _ := gval.Evaluate("foo1.xyz > 5 && bar == \"baz\"", params, lang) // ✅ false
result, _ := gval.Evaluate("foo1.xyz > 5 || bar == \"baz\"", params, lang) // ✅ true
result, _ := gval.Evaluate("foo1.xyz > 5 || bar == \"different\"", params, lang) // ✅ false
```

## Behavior

### Missing Field Cases:
- `missing_field > 5` → `false` (nil > 5 = false)
- `missing_field > 5 && bar == "baz"` → `false` (false && true = false)
- `missing_field > 5 || bar == "baz"` → `true` (false || true = true)  
- `missing_field > 5 || bar == "different"` → `false` (false || false = false)

### Existing Field Cases (unchanged):
- `existing_field > 5 && bar == "baz"` → works normally based on actual values
- `existing_field > 5 || bar == "baz"` → works normally based on actual values

## Alternative Approach

If you need more granular control over logical operators, you can also use `MissingFieldTolerantLogic()` which directly handles missing field errors in `&&` and `||` operators:

```go
lang := gval.Full(gval.MissingFieldTolerantLogic())

// This approach intercepts "unknown parameter" errors in logical expressions
// and treats them as false, allowing the logical operation to continue
```

## Files Added

- `missing_field_tolerant.go` - Contains the new language extensions
- `missing_field_test.go` - Comprehensive tests for the new functionality  
- `demo_missing_fields.go` - Example usage demonstration

## Integration

These extensions are backward compatible and don't affect existing functionality. They only change behavior for missing field access, making expressions more robust and predictable.
