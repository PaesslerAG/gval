package gval

import "testing"

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
