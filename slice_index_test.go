package gval

import (
	"reflect"
	"strings"
	"testing"
)

func TestSliceIndexSelection(t *testing.T) {
	generic := map[string]interface{}{"items": []interface{}{"first", "last"}}
	typed := map[string]interface{}{"items": []string{"first", "last"}}

	tests := []struct {
		name       string
		expression string
		parameter  interface{}
		want       interface{}
		wantErr    bool
	}{
		{name: "generic negative", expression: "items[-1]", parameter: generic, wantErr: true},
		{name: "generic out of range", expression: "items[2]", parameter: generic, wantErr: true},
		{name: "generic noninteger", expression: `items["bad"]`, parameter: generic, wantErr: true},
		{name: "generic fractional", expression: "items[1.5]", parameter: generic, wantErr: true},
		{name: "generic nested invalid index", expression: "items[99][0]", parameter: generic, wantErr: true},
		{name: "generic empty", expression: "items[0]", parameter: map[string]interface{}{"items": []interface{}{}}, wantErr: true},
		{name: "generic nil slice", expression: "items[0]", parameter: map[string]interface{}{"items": []interface{}(nil)}, wantErr: true},
		{name: "typed negative", expression: "items[-1]", parameter: typed, wantErr: true},
		{name: "typed out of range", expression: "items[2]", parameter: typed, wantErr: true},
		{name: "typed noninteger", expression: `items["bad"]`, parameter: typed, wantErr: true},
		{name: "typed nil slice", expression: "items[0]", parameter: map[string]interface{}{"items": []string(nil)}, wantErr: true},
		{name: "generic first", expression: "items[0]", parameter: generic, want: "first"},
		{name: "generic last", expression: "items[1]", parameter: generic, want: "last"},
		{name: "typed first", expression: "items[0]", parameter: typed, want: "first"},
		{name: "typed last", expression: "items[1]", parameter: typed, want: "last"},
		{name: "valid nil element", expression: "items[0]", parameter: map[string]interface{}{"items": []interface{}{nil}}},
		{name: "missing string map key", expression: "missing", parameter: map[string]interface{}{}},
		{name: "missing interface map key", expression: "missing", parameter: map[interface{}]interface{}{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.expression, tt.parameter)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Evaluate(%q) = %#v, nil; want unknown parameter error", tt.expression, got)
				}
				if !strings.Contains(err.Error(), "unknown parameter") {
					t.Fatalf("Evaluate(%q) error = %v; want unknown parameter error", tt.expression, err)
				}
				if got != nil {
					t.Fatalf("Evaluate(%q) = %#v with error; want nil", tt.expression, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Evaluate(%q) error = %v", tt.expression, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Evaluate(%q) = %#v; want %#v", tt.expression, got, tt.want)
			}
		})
	}
}
