package gval

import (
	"testing"
)

func TestUserRequestedBehavior(t *testing.T) {
	// Your exact test case
	params := map[string]interface{}{
		"foo":  10,
		"bar":  "baz",
		"foo1": map[string]interface{}{"xyz1": 100},
	}
	
	// Create the tolerant language
	lang := TolerantFull()

	// Test Case 1: foo1.xyz > 5 && bar == baz should return false
	t.Run("Case1: foo1.xyz > 5 && bar == baz should return false", func(t *testing.T) {
		result, err := Evaluate("foo1.xyz > 5 && bar == \"baz\"", params, lang)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result != false {
			t.Errorf("Expected false, got %v", result)
		}
	})

	// Test Case 3: foo1.xyz > 5 || bar == baz should return true  
	t.Run("Case3: foo1.xyz > 5 || bar == baz should return true", func(t *testing.T) {
		result, err := Evaluate("foo1.xyz > 5 || bar == \"baz\"", params, lang)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result != true {
			t.Errorf("Expected true, got %v", result)
		}
	})

	// Additional test: both conditions false
	t.Run("Both conditions false should return false", func(t *testing.T) {
		result, err := Evaluate("foo1.xyz > 5 || bar == \"different\"", params, lang)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result != false {
			t.Errorf("Expected false, got %v", result)
		}
	})

	// Test that existing fields still work
	t.Run("Existing fields should work normally", func(t *testing.T) {
		result, err := Evaluate("foo1.xyz1 > 50 && bar == \"baz\"", params, lang)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result != true {
			t.Errorf("Expected true, got %v", result)
		}
	})

	// Test direct missing field access
	t.Run("Direct missing field access should return false", func(t *testing.T) {
		result, err := Evaluate("foo1.xyz", params, lang)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result != false {
			t.Errorf("Expected false, got %v", result)
		}
	})
}
