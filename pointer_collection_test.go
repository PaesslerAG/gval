package gval_test

import (
	"testing"

	"github.com/PaesslerAG/gval"
)

func TestPointerCollectionSelection(t *testing.T) {
	mapping := map[string]int{"value": 42}
	slice := []int{42}
	for _, tc := range []struct {
		name       string
		expression string
		input      interface{}
	}{
		{"map pointer", "data.value", &mapping},
		{"slice pointer", "data[0]", &slice},
		{"map value", "data.value", mapping},
		{"slice value", "data[0]", slice},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := gval.Evaluate(tc.expression, map[string]interface{}{"data": tc.input})
			if err != nil || got != 42 {
				t.Fatalf("got %v, %v; want 42", got, err)
			}
		})
	}
}
