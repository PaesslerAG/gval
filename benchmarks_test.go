package gval

import (
	"context"
	"testing"
)

func BenchmarkGval(bench *testing.B) {
	benchmarks := []evaluationTest{
		{
			// Serves as a "water test" to give an idea of the general overhead
			name:       "const",
			expression: "1",
		},
		{
			name:       "single parameter",
			expression: "requests_made",
			parameter: map[string]interface{}{
				"requests_made": 99.0,
			},
		},
		{
			name:       "parameter",
			expression: "requests_made > requests_succeeded",
			parameter: map[string]interface{}{
				"requests_made":      99.0,
				"requests_succeeded": 90.0,
			},
		},
		{
			// The most common use case, a single variable, modified slightly, compared to a constant.
			// This is the "expected" use case.
			name:       "common",
			expression: "(requests_made * requests_succeeded / 100) >= 90",
			parameter: map[string]interface{}{
				"requests_made":      99.0,
				"requests_succeeded": 90.0,
			},
		},
		{
			// All major possibilities in one expression.
			name: "complex",
			expression: `2 > 1 &&
			"something" != "nothing" ||
			date("2014-01-20") < date("Wed Jul  8 23:07:35 MDT 2015") && 
			object["Variable name with spaces"] <= array[0] &&
			modifierTest + 1000 / 2 > (80 * 100 % 2)`,
			parameter: map[string]interface{}{
				"object":       map[string]interface{}{"Variable name with spaces": 10.},
				"array":        []interface{}{0.},
				"modifierTest": 7.3,
			},
		},
		{
			// no variables, no modifiers
			name:       "literal",
			expression: "(2) > (1)",
		},
		{
			name:       "modifier",
			expression: "(2) + (2) == (4)",
		},
		{
			//   Benchmarks uncompiled parameter regex operators, which are the most expensive of the lot.
			//   Note that regex compilation times are unpredictable and wily things. The regex engine has a lot of edge cases
			//   and possible performance pitfalls. This test doesn't aim to be comprehensive against all possible regex scenarios,
			//   it is primarily concerned with tracking how much longer it takes to compile a regex at evaluation-time than during parse-time.
			name:       "regex",
			expression: "(foo !~ bar) && (foo + bar =~ oba)",
			parameter: map[string]interface{}{
				"foo": "foo",
				"bar": "bar",
				"baz": "baz",
				"oba": ".*oba.*",
			},
		},
		{
			// Benchmarks pre-compilable regex patterns. Meant to serve as a sanity check that constant strings used as regex patterns
			// are actually being precompiled.
			// Also demonstrates that (generally) compiling a regex at evaluation-time takes an order of magnitude more time than pre-compiling.
			name:       "constant regex",
			expression: `(foo !~ "[bB]az") && (bar =~ "[bB]ar")`,
			parameter: map[string]interface{}{
				"foo": "foo",
				"bar": "bar",
				"baz": "baz",
				"oba": ".*oba.*",
			},
		},
		{
			name:       "accessors",
			expression: "foo.Int",
			parameter:  fooFailureParameters,
		},
		{
			name:       "accessors method",
			expression: "foo.Func()",
			parameter:  fooFailureParameters,
		},
		{
			name:       "accessors method parameter",
			expression: `foo.FuncArgStr("bonk")`,
			parameter:  fooFailureParameters,
		},
		{
			name:       "nested accessors",
			expression: `foo.Nested.Funk`,
			parameter:  fooFailureParameters,
		},
	}
	for _, benchmark := range benchmarks {
		eval, err := Full().NewEvaluable(benchmark.expression)
		if err != nil {
			bench.Fatal(err)
		}
		_, err = eval(context.Background(), benchmark.parameter)
		if err != nil {
			bench.Fatal(err)
		}
		bench.Run(benchmark.name+"_evaluation", func(bench *testing.B) {
			for i := 0; i < bench.N; i++ {
				eval(context.Background(), benchmark.parameter)
			}
		})
		bench.Run(benchmark.name+"_parsing", func(bench *testing.B) {
			for i := 0; i < bench.N; i++ {
				Full().NewEvaluable(benchmark.expression)
			}
		})

	}
}
