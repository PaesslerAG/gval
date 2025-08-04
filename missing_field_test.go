package gval

import (
	"reflect"
	"strings"
	"testing"
)

func TestMissingFieldHandling(t *testing.T) {
	params := map[string]interface{}{
		"foo":  10,
		"bar":  "baz",
		"foo1": map[string]interface{}{"xyz1": 100},
	}

	// Test the default behavior (should fail)
	t.Run("Default behavior (should fail)", func(t *testing.T) {
		_, err := Evaluate("foo1.xyz > 5", params)
		if err == nil {
			t.Error("Expected error for missing field but got none")
		}
		// Accept either "unknown parameter" or "invalid operation" errors
		errStr := err.Error()
		if !strings.Contains(errStr, "unknown parameter") && !strings.Contains(errStr, "invalid operation") {
			t.Errorf("Expected 'unknown parameter' or 'invalid operation' error but got: %v", err)
		}
	})

	// Test with TolerantFull approach - the recommended solution
	t.Run("TolerantFull approach", func(t *testing.T) {
		lang := TolerantFull()
		
		tests := []struct {
			name       string
			expression string
			want       interface{}
		}{
			{
				name:       "Missing field comparison should return false",
				expression: "foo1.xyz > 5",
				want:       false,
			},
			{
				name:       "Missing field with AND - should return false",
				expression: "foo1.xyz > 5 && bar == \"baz\"",
				want:       false,
			},
			{
				name:       "Missing field with OR - should return true",
				expression: "foo1.xyz > 5 || bar == \"baz\"",
				want:       true,
			},
			{
				name:       "Missing field with OR (both false) - should return false",
				expression: "foo1.xyz > 5 || bar == \"different\"",
				want:       false,
			},
			{
				name:       "Existing field should work normally",
				expression: "foo1.xyz1 > 50 && bar == \"baz\"",
				want:       true,
			},
			{
				name:       "Missing top-level field with AND",
				expression: "missing > 5 && bar == \"baz\"",
				want:       false,
			},
			{
				name:       "Missing top-level field with OR",
				expression: "missing > 5 || bar == \"baz\"",
				want:       true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := Evaluate(tt.expression, params, lang)
				if err != nil {
					t.Errorf("Evaluate() error = %v", err)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Evaluate(%s) = %v, want %v", tt.expression, got, tt.want)
				}
			})
		}
	})

	// Test with specific missing field behaviors
	t.Run("Different missing field behaviors", func(t *testing.T) {
		tests := []struct {
			name      string
			behavior  Language
			wantValue interface{}
			wantError bool
		}{
			{
				name:      "FalseOnMissingField",
				behavior:  Full(WithMissingFieldBehavior(FalseOnMissingField)),
				wantValue: false,
				wantError: false,
			},
			{
				name:      "NilOnMissingField", 
				behavior:  Full(WithMissingFieldBehavior(NilOnMissingField)),
				wantValue: nil,
				wantError: false,
			},
			{
				name:      "ErrorOnMissingField (default)",
				behavior:  Full(WithMissingFieldBehavior(ErrorOnMissingField)),
				wantValue: nil,
				wantError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := Evaluate("foo1.xyz", params, tt.behavior)
				if tt.wantError {
					if err == nil {
						t.Error("Expected error but got none")
					}
				} else {
					if err != nil {
						t.Errorf("Unexpected error: %v", err)
					}
					if !reflect.DeepEqual(got, tt.wantValue) {
						t.Errorf("Got %v, want %v", got, tt.wantValue)
					}
				}
			})
		}
	})
}
