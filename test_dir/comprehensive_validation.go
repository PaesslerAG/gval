package main

import (
	"fmt"
	"time"
	"github.com/Nandagopi/gval"
)

func main() {
	params := map[string]interface{}{
		"foo":  10,
		"bar":  "baz",
		"foo1": map[string]interface{}{"xyz1": 100},
	}

	fmt.Println("=== COMPREHENSIVE VALIDATION ===")
	
	// Test your exact use cases
	lang := gval.TolerantFull()
	
	testCases := []struct {
		name     string
		expr     string
		expected interface{}
	}{
		{"Missing field comparison", "foo1.xyz > 5", false},
		{"Case 1: AND with missing field", "foo1.xyz > 5 && bar == \"baz\"", false},
		{"Case 3: OR with missing field", "foo1.xyz > 5 || bar == \"baz\"", true},
		{"Both conditions false", "foo1.xyz > 5 || bar == \"different\"", false},
		{"Existing field works", "foo1.xyz1 > 50 && bar == \"baz\"", true},
		{"Complex expression", "(foo1.xyz > 5 || foo > 8) && bar == \"baz\"", true},
		{"Nested missing", "foo1.missing.deep > 5", false},
		{"Direct missing field", "foo1.xyz", false},
	}

	allPassed := true
	for _, tc := range testCases {
		result, err := gval.Evaluate(tc.expr, params, lang)
		if err != nil {
			fmt.Printf("❌ %s: ERROR - %v\n", tc.name, err)
			allPassed = false
		} else if result != tc.expected {
			fmt.Printf("❌ %s: FAIL - got %v, expected %v\n", tc.name, result, tc.expected)
			allPassed = false
		} else {
			fmt.Printf("✅ %s: PASS - %v\n", tc.name, result)
		}
	}

	fmt.Printf("\n=== OVERALL RESULT: %s ===\n", map[bool]string{true: "ALL TESTS PASSED", false: "SOME TESTS FAILED"}[allPassed])

	// Performance comparison
	fmt.Println("\n=== PERFORMANCE ANALYSIS ===")
	
	expr := "foo1.xyz > 5 || bar == \"baz\""
	iterations := 10000
	
	// Test TolerantFull performance
	start := time.Now()
	for i := 0; i < iterations; i++ {
		gval.Evaluate(expr, params, lang)
	}
	tolerantDuration := time.Since(start)
	
	// Test standard Full performance (with error)
	standardLang := gval.Full()
	start = time.Now()
	for i := 0; i < iterations; i++ {
		gval.Evaluate("foo1.xyz1 > 50 || bar == \"baz\"", params, standardLang) // Use existing field
	}
	standardDuration := time.Since(start)
	
	fmt.Printf("TolerantFull (with missing fields): %v (%v per operation)\n", 
		tolerantDuration, time.Duration(int64(tolerantDuration)/int64(iterations)))
	fmt.Printf("Standard Full (existing fields): %v (%v per operation)\n", 
		standardDuration, time.Duration(int64(standardDuration)/int64(iterations)))
	
	overhead := float64(tolerantDuration) / float64(standardDuration)
	fmt.Printf("Performance overhead: %.2fx\n", overhead)
	
	if overhead < 1.5 {
		fmt.Println("✅ Performance impact is acceptable (< 1.5x)")
	} else {
		fmt.Println("⚠️  Performance impact is significant (> 1.5x)")
	}
}
