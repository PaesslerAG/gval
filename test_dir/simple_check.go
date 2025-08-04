package main

import (
	"fmt"
	"strings"
)

// Import only the standard gval to test default behavior
import "github.com/Nandagopi/gval"

func main() {
	params := map[string]interface{}{
		"foo1": map[string]interface{}{"xyz1": 100},
	}
	
	// Test just the missing field access
	_, err := gval.Evaluate("foo1.xyz", params)
	fmt.Printf("Error: %v\n", err)
	fmt.Printf("Contains 'unknown parameter': %v\n", strings.Contains(err.Error(), "unknown parameter"))
	
	// Test the comparison that's failing
	_, err2 := gval.Evaluate("foo1.xyz > 5", params) 
	fmt.Printf("Comparison Error: %v\n", err2)
}
