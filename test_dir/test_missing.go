package main

import (
	"fmt"
	"github.com/Nandagopi/gval"
)

func main() {
	params := map[string]interface{}{
		"foo": 10,
		"bar": "baz",
		"foo1": map[string]interface{}{"xyz1": 100},
	}
	
	// Test case 1: Missing field (foo1.xyz)
	fmt.Println("=== Test Case 1: foo1.xyz > 5 ===")
	val, err := gval.Evaluate("foo1.xyz > 5", params)
	fmt.Printf("Result: %v, Error: %v\n\n", val, err)
	
	// Test case 2: Missing field with AND
	fmt.Println("=== Test Case 2: foo1.xyz > 5 && bar == \"baz\" ===")
	val, err = gval.Evaluate("foo1.xyz > 5 && bar == \"baz\"", params)
	fmt.Printf("Result: %v, Error: %v\n\n", val, err)
	
	// Test case 3: Missing field with OR
	fmt.Println("=== Test Case 3: foo1.xyz > 5 || bar == \"baz\" ===")
	val, err = gval.Evaluate("foo1.xyz > 5 || bar == \"baz\"", params)
	fmt.Printf("Result: %v, Error: %v\n\n", val, err)
	
	// Test case 4: Existing field
	fmt.Println("=== Test Case 4: foo1.xyz1 > 5 ===")
	val, err = gval.Evaluate("foo1.xyz1 > 5", params)
	fmt.Printf("Result: %v, Error: %v\n\n", val, err)
}
