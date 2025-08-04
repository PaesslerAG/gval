package main

import (
	"fmt"
	"github.com/Nandagopi/gval"
)

func main() {
	params := map[string]interface{}{
		"foo":  10,
		"bar":  "baz",
		"foo1": map[string]interface{}{"xyz1": 100},
	}
	
	// Test default behavior
	fmt.Println("=== Default Behavior ===")
	result, err := gval.Evaluate("foo1.xyz", params)
	fmt.Printf("foo1.xyz: Result = %v, Error = %v\n", result, err)
	
	result, err = gval.Evaluate("foo1.xyz > 5", params)
	fmt.Printf("foo1.xyz > 5: Result = %v, Error = %v\n", result, err)
}
