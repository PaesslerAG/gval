package gval

import (
	"context"
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
