package gval

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestEvaluable_IsConst(t *testing.T) {
	p := Parser{}
	tests := []struct {
		name string
		e    Evaluable
		want bool
	}{
		{
			"const",
			p.Const(80.5),
			true,
		},
		{
			"var",
			p.Var(),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsConst(); got != tt.want {
				t.Errorf("Evaluable.IsConst() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvaluable_EvalInt(t *testing.T) {
	tests := []struct {
		name    string
		e       Evaluable
		want    int
		wantErr bool
	}{
		{
			"point",
			constant("5.3"),
			5,
			false,
		},
		{
			"number",
			constant(255.),
			255,
			false,
		},
		{
			"error",
			constant("5.3 cm"),
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.EvalInt(context.Background(), nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluable.EvalInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Evaluable.EvalInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvaluable_EvalFloat64(t *testing.T) {
	tests := []struct {
		name    string
		e       Evaluable
		want    float64
		wantErr bool
	}{
		{
			"point",
			constant("5.3"),
			5.3,
			false,
		},
		{
			"number",
			constant(255.),
			255,
			false,
		},
		{
			"error",
			constant("5.3 cm"),
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.EvalFloat64(context.Background(), nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluable.EvalFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Evaluable.EvalFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

type custSel struct {
	Str string
	Map map[string]interface{}
}

func (s custSel) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	if k == "str" {
		return s.Str, nil
	}

	if k == "map" {
		return s.Map, nil
	}

	if strings.HasPrefix(k, "deep") {
		return s, nil
	}

	return nil, fmt.Errorf("unknown-key")

}

func TestEvaluable_CustomSelector(t *testing.T) {
	var (
		lang  = Base()
		tests = []struct {
			name    string
			expr    string
			params  interface{}
			want    interface{}
			wantErr bool
		}{
			{
				"unknown",
				"s.Foo",
				map[string]interface{}{"s": &custSel{}},
				nil,
				true,
			},
			{
				"field directly",
				"s.Str",
				map[string]interface{}{"s": &custSel{Str: "test-value"}},
				nil,
				true,
			},
			{
				"field via selector",
				"s.str",
				map[string]interface{}{"s": &custSel{Str: "test-value"}},
				"test-value",
				false,
			},
			{
				"flat",
				"str",
				&custSel{Str: "test-value"},
				"test-value",
				false,
			},
			{
				"map field",
				"s.map.foo",
				map[string]interface{}{"s": &custSel{Map: map[string]interface{}{"foo": "bar"}}},
				"bar",
				false,
			},
			{
				"crawl to val",
				"deep.deeper.deepest.str",
				&custSel{Str: "foo"},
				"foo",
				false,
			},
			{
				"crawl to struct",
				"deep.deeper.deepest",
				&custSel{},
				custSel{},
				false,
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lang.Evaluate(tt.expr, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluable.Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Evaluable.Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
