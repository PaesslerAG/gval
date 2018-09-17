package gval

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

func Test_Infix(t *testing.T) {
	type subTest struct {
		name    string
		a       interface{}
		b       interface{}
		wantRet interface{}
	}
	tests := []struct {
		name string
		infix
		subTests []subTest
	}{
		{
			"number operator",
			infix{
				number: func(a, b float64) (interface{}, error) { return a * b, nil },
			},
			[]subTest{
				{"float64 arguments", 7., 3., 21.},
				{"int arguments", 7, 3, 21.},
				{"string arguments", "7", "3.", 21.},
			},
		},
		{
			"number and string operator",
			infix{
				number: func(a, b float64) (interface{}, error) { return a + b, nil },
				text:   func(a, b string) (interface{}, error) { return fmt.Sprintf("%v%v", a, b), nil },
			},

			[]subTest{
				{"float64 arguments", 7., 3., 10.},
				{"int arguments", 7, 3, 10.},
				{"number string arguments", "7", "3.", "73."},
				{"string arguments", "hello ", "world", "hello world"},
			},
		},
		{
			"bool operator",
			infix{
				shortCircuit: func(a interface{}) (interface{}, bool) { return false, a == false },
				boolean:      func(a, b bool) (interface{}, error) { return a && b, nil },
			},

			[]subTest{
				{"bool arguments", false, true, false},
				{"number arguments", 0, true, false},
				{"lower string arguments", "false", "true", false},
				{"upper string arguments", "TRUE", "FALSE", false},
				{"shortCircuit", false, "not a boolean", false},
			},
		},
		{
			"bool, number, text and interface operator",
			infix{
				number:    func(a, b float64) (interface{}, error) { return a == b, nil },
				boolean:   func(a, b bool) (interface{}, error) { return a == b, nil },
				text:      func(a, b string) (interface{}, error) { return a == b, nil },
				arbitrary: func(a, b interface{}) (interface{}, error) { return a == b, nil },
			},

			[]subTest{
				{"number string and int arguments", "7", 7, true},
				{"bool string and bool arguments", "true", true, true},
				{"string arguments", "hello", "hello", true},
				{"upper string arguments", "TRUE", "FALSE", false},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.infix.initiate("<" + tt.name + ">")
			builder := tt.infix.builder
			for _, tt := range tt.subTests {
				t.Run(tt.name, func(t *testing.T) {
					eval, err := builder(constant(tt.a), constant(tt.b))
					if err != nil {
						t.Fatal(err)
					}

					got, err := eval(context.Background(), nil)
					if err != nil {
						t.Fatal(err)
					}

					if !reflect.DeepEqual(got, tt.wantRet) {
						t.Fatalf("binaryOperator() eval() = %v, want %v", got, tt.wantRet)
					}
				})
			}
		})
	}
}

func Test_stageStack_push(t *testing.T) {
	p := (*Parser)(nil)
	tests := []struct {
		name   string
		pres   []operatorPrecedence
		expect string
	}{
		{
			"flat",
			[]operatorPrecedence{1, 1, 1, 1},
			"((((AB)C)D)E)",
		},
		{
			"asc",
			[]operatorPrecedence{1, 2, 3, 4},
			"(A(B(C(DE))))",
		},
		{
			"desc",
			[]operatorPrecedence{4, 3, 2, 1},
			"((((AB)C)D)E)",
		},
		{
			"mixed",
			[]operatorPrecedence{1, 2, 1, 1},
			"(((A(BC))D)E)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			X := int('A')

			op := func(a, b Evaluable) (Evaluable, error) {
				return func(c context.Context, o interface{}) (interface{}, error) {
					aa, _ := a.EvalString(c, nil)
					bb, _ := b.EvalString(c, nil)
					s := "(" + aa + bb + ")"
					return s, nil
				}, nil
			}
			stack := stageStack{}
			for _, pre := range tt.pres {
				if err := stack.push(stage{p.Const(string(rune(X))), op, pre}); err != nil {
					t.Fatal(err)
				}
				X++
			}

			if err := stack.push(stage{p.Const(string(rune(X))), nil, 0}); err != nil {
				t.Fatal(err)
			}

			if len(stack) != 1 {
				t.Fatalf("stack must hold exactly one element")
			}

			got, _ := stack[0].EvalString(context.Background(), nil)
			if got != tt.expect {
				t.Fatalf("got %s but expected %s", got, tt.expect)
			}
		})
	}
}
