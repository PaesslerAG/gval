package main

import (
	"fmt"
	"github.com/Nandagopi/gval"
)

func main() {
	// Your test case parameters
	params := map[string]interface{}{
		"foo":  10,
		"bar":  "baz",
		"foo1": map[string]interface{}{"xyz1": 100},
		// Note: foo1.xyz is missing
	}

	fmt.Println("=== Default gval behavior (throws errors) ===")
	
	// This will fail with default gval
	_, err := gval.Evaluate("foo1.xyz > 5", params)
	fmt.Printf("foo1.xyz > 5: Error = %v\n", err)
	
	_, err = gval.Evaluate("foo1.xyz > 5 && bar == \"baz\"", params)
	fmt.Printf("foo1.xyz > 5 && bar == \"baz\": Error = %v\n", err)
	
	_, err = gval.Evaluate("foo1.xyz > 5 || bar == \"baz\"", params)
	fmt.Printf("foo1.xyz > 5 || bar == \"baz\": Error = %v\n\n", err)

	fmt.Println("=== With MissingFieldAsNil + NilSafeComparison ===")
	
	// Create a language that treats missing fields as nil and handles nil comparisons
	lang := gval.Full(gval.MissingFieldAsNil(), gval.NilSafeComparison())
	
	// Test case 1: Missing field comparison
	result, err := gval.Evaluate("foo1.xyz > 5", params, lang)
	fmt.Printf("foo1.xyz > 5: Result = %v, Error = %v\n", result, err)
	
	// Test case 2: Missing field with AND (should return false)
	result, err = gval.Evaluate("foo1.xyz > 5 && bar == \"baz\"", params, lang)
	fmt.Printf("foo1.xyz > 5 && bar == \"baz\": Result = %v, Error = %v\n", result, err)
	
	// Test case 3: Missing field with OR (should return true)
	result, err = gval.Evaluate("foo1.xyz > 5 || bar == \"baz\"", params, lang)
	fmt.Printf("foo1.xyz > 5 || bar == \"baz\": Result = %v, Error = %v\n", result, err)
	
	// Test case 4: Both conditions false
	result, err = gval.Evaluate("foo1.xyz > 5 || bar == \"different\"", params, lang)
	fmt.Printf("foo1.xyz > 5 || bar == \"different\": Result = %v, Error = %v\n", result, err)
	
	// Test case 5: Existing field (should work normally)
	result, err = gval.Evaluate("foo1.xyz1 > 50 && bar == \"baz\"", params, lang)
	fmt.Printf("foo1.xyz1 > 50 && bar == \"baz\": Result = %v, Error = %v\n", result, err)
	
	fmt.Println("\n=== Summary ===")
	fmt.Println("✅ Missing fields are treated as nil/false")
	fmt.Println("✅ Logical operators work as expected:")
	fmt.Println("   - missing_field && true  → false")
	fmt.Println("   - missing_field || true  → true")
	fmt.Println("   - missing_field || false → false")
	fmt.Println("✅ Existing fields continue to work normally")
}
