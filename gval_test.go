package gval

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type evaluationTest struct {
	name       string
	expression string
	extension  Language
	parameter  interface{}
	want       interface{}
	wantErr    string
}

func testEvaluate(tests []evaluationTest, t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.expression, tt.parameter, tt.extension)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("Evaluate(%s) expected error but got %v", tt.expression, got)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("Evaluate(%s) expected error %s but got error %v", tt.expression, tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("Evaluate() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Evaluate(%s) = %v, want %v", tt.expression, got, tt.want)
			}
		})
	}
}

//dummyParameter used to test "parameter calls".
type dummyParameter struct {
	String    string
	Int       int
	BoolFalse bool
	Nil       interface{}
	Nested    dummyNestedParameter
}

func (d dummyParameter) Func() string {
	return "funk"
}

func (d dummyParameter) Func2() (string, error) {
	return "frink", nil
}

func (d *dummyParameter) PointerFunc() (string, error) {
	return "point", nil
}

func (d dummyParameter) FuncErr() (string, error) {
	return "", fmt.Errorf("fumps")
}

func (d dummyParameter) FuncArgStr(arg1 string) string {
	return arg1
}

func (d dummyParameter) AlwaysFail() (interface{}, error) {
	return nil, fmt.Errorf("function should always fail")
}

type dummyNestedParameter struct {
	Funk  string
	Map   map[string]int
	Slice []int
}

func (d dummyNestedParameter) Dunk(arg1 string) string {
	return arg1 + "dunk"
}

var foo = dummyParameter{
	String:    "string!",
	Int:       101,
	BoolFalse: false,
	Nil:       nil,
	Nested: dummyNestedParameter{
		Funk:  "funkalicious",
		Map:   map[string]int{"a": 1, "b": 2, "c": 3},
		Slice: []int{1, 2, 3},
	},
}

var fooFailureParameters = map[string]interface{}{
	"foo":    foo,
	"fooptr": &foo,
}
