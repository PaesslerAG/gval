package gval

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestParameterized(t *testing.T) {
	testEvaluate(
		[]evaluationTest{
			{
				name:       "Single parameter modified by constant",
				expression: "foo + 2",
				parameter: map[string]interface{}{
					"foo": 2.0,
				},
				want: 4.0,
			},
			{

				name:       "Single parameter modified by variable",
				expression: "foo * bar",
				parameter: map[string]interface{}{
					"foo": 5.0,
					"bar": 2.0,
				},
				want: 10.0,
			},
			{

				name:       "Single parameter modified by variable",
				expression: `foo["hey"] * bar[1]`,
				parameter: map[string]interface{}{
					"foo": map[string]interface{}{"hey": 5.0},
					"bar": []interface{}{7., 2.0},
				},
				want: 10.0,
			},
			{

				name:       "Multiple multiplications of the same parameter",
				expression: "foo * foo * foo",
				parameter: map[string]interface{}{
					"foo": 10.0,
				},
				want: 1000.0,
			},
			{

				name:       "Multiple additions of the same parameter",
				expression: "foo + foo + foo",
				parameter: map[string]interface{}{
					"foo": 10.0,
				},
				want: 30.0,
			},
			{

				name:       "Parameter name sensitivity",
				expression: "foo + FoO + FOO",
				parameter: map[string]interface{}{
					"foo": 8.0,
					"FoO": 4.0,
					"FOO": 2.0,
				},
				want: 14.0,
			},
			{

				name:       "Sign prefix comparison against prefixed variable",
				expression: "-1 < -foo",
				parameter:  map[string]interface{}{"foo": -8.0},
				want:       true,
			},
			{

				name:       "Fixed-point parameter",
				expression: "foo > 1",
				parameter:  map[string]interface{}{"foo": 2},
				want:       true,
			},
			{

				name:       "Modifier after closing clause",
				expression: "(2 + 2) + 2 == 6",
				want:       true,
			},
			{

				name:       "Comparator after closing clause",
				expression: "(2 + 2) >= 4",
				want:       true,
			},
			{

				name:       "Two-boolean logical operation (for issue #8)",
				expression: "(foo == true) || (bar == true)",
				parameter: map[string]interface{}{
					"foo": true,
					"bar": false,
				},
				want: true,
			},
			{

				name:       "Two-variable integer logical operation (for issue #8)",
				expression: "foo > 10 && bar > 10",
				parameter: map[string]interface{}{
					"foo": 1,
					"bar": 11,
				},
				want: false,
			},
			{

				name:       "Regex against right-hand parameter",
				expression: `"foobar" =~ foo`,
				parameter: map[string]interface{}{
					"foo": "obar",
				},
				want: true,
			},
			{

				name:       "Not-regex against right-hand parameter",
				expression: `"foobar" !~ foo`,
				parameter: map[string]interface{}{
					"foo": "baz",
				},
				want: true,
			},
			{

				name:       "Regex against two parameter",
				expression: `foo =~ bar`,
				parameter: map[string]interface{}{
					"foo": "foobar",
					"bar": "oba",
				},
				want: true,
			},
			{

				name:       "Not-regex against two parameter",
				expression: "foo !~ bar",
				parameter: map[string]interface{}{
					"foo": "foobar",
					"bar": "baz",
				},
				want: true,
			},
			{

				name:       "Pre-compiled regex",
				expression: "foo =~ bar",
				parameter: map[string]interface{}{
					"foo": "foobar",
					"bar": regexp.MustCompile("[fF][oO]+"),
				},
				want: true,
			},
			{

				name:       "Pre-compiled not-regex",
				expression: "foo !~ bar",
				parameter: map[string]interface{}{
					"foo": "foobar",
					"bar": regexp.MustCompile("[fF][oO]+"),
				},
				want: false,
			},
			{

				name:       "Single boolean parameter",
				expression: "commission ? 10",
				parameter: map[string]interface{}{
					"commission": true},
				want: 10.0,
			},
			{

				name:       "True comparator with a parameter",
				expression: `partner == "amazon" ? 10`,
				parameter: map[string]interface{}{
					"partner": "amazon"},
				want: 10.0,
			},
			{

				name:       "False comparator with a parameter",
				expression: `partner == "amazon" ? 10`,
				parameter: map[string]interface{}{
					"partner": "ebay"},
				want: nil,
			},
			{

				name:       "True comparator with multiple parameters",
				expression: "theft && period == 24 ? 60",
				parameter: map[string]interface{}{
					"theft":  true,
					"period": 24,
				},
				want: 60.0,
			},
			{

				name:       "False comparator with multiple parameters",
				expression: "theft && period == 24 ? 60",
				parameter: map[string]interface{}{
					"theft":  false,
					"period": 24,
				},
				want: nil,
			},
			{

				name:       "String concat with single string parameter",
				expression: `foo + "bar"`,
				parameter: map[string]interface{}{
					"foo": "baz"},
				want: "bazbar",
			},
			{

				name:       "String concat with multiple string parameter",
				expression: "foo + bar",
				parameter: map[string]interface{}{
					"foo": "baz",
					"bar": "quux",
				},
				want: "bazquux",
			},
			{

				name:       "String concat with float parameter",
				expression: "foo + bar",
				parameter: map[string]interface{}{
					"foo": "baz",
					"bar": 123.0,
				},
				want: "baz123",
			},
			{

				name:       "Mixed multiple string concat",
				expression: `foo + 123 + "bar" + true`,
				parameter:  map[string]interface{}{"foo": "baz"},
				want:       "baz123bartrue",
			},
			{

				name:       "Integer width spectrum",
				expression: "uint8 + uint16 + uint32 + uint64 + int8 + int16 + int32 + int64",
				parameter: map[string]interface{}{
					"uint8":  uint8(0),
					"uint16": uint16(0),
					"uint32": uint32(0),
					"uint64": uint64(0),
					"int8":   int8(0),
					"int16":  int16(0),
					"int32":  int32(0),
					"int64":  int64(0),
				},
				want: 0.0,
			},
			{

				name:       "Null coalesce right",
				expression: "foo ?? 1.0",
				parameter:  map[string]interface{}{"foo": nil},
				want:       1.0,
			},
			{

				name:       "Multiple comparator/logical operators (#30)",
				expression: "(foo >= 2887057408 && foo <= 2887122943) || (foo >= 168100864 && foo <= 168118271)",
				parameter:  map[string]interface{}{"foo": 2887057409},
				want:       true,
			},
			{

				name:       "Multiple comparator/logical operators, opposite order (#30)",
				expression: "(foo >= 168100864 && foo <= 168118271) || (foo >= 2887057408 && foo <= 2887122943)",
				parameter:  map[string]interface{}{"foo": 2887057409},
				want:       true,
			},
			{

				name:       "Multiple comparator/logical operators, small value (#30)",
				expression: "(foo >= 2887057408 && foo <= 2887122943) || (foo >= 168100864 && foo <= 168118271)",
				parameter:  map[string]interface{}{"foo": 168100865},
				want:       true,
			},
			{

				name:       "Multiple comparator/logical operators, small value, opposite order (#30)",
				expression: "(foo >= 168100864 && foo <= 168118271) || (foo >= 2887057408 && foo <= 2887122943)",
				parameter:  map[string]interface{}{"foo": 168100865},
				want:       true,
			},
			{

				name:       "Incomparable array equality comparison",
				expression: "arr == arr",
				parameter:  map[string]interface{}{"arr": []int{0, 0, 0}},
				want:       true,
			},
			{

				name:       "Incomparable array not-equality comparison",
				expression: "arr != arr",
				parameter:  map[string]interface{}{"arr": []int{0, 0, 0}},
				want:       false,
			},
			{

				name:       "Mixed function and parameters",
				expression: "sum(1.2, amount) + name",
				extension: Function("sum", func(arguments ...interface{}) (interface{}, error) {
					sum := 0.0
					for _, v := range arguments {
						sum += v.(float64)
					}
					return sum, nil
				},
				),
				parameter: map[string]interface{}{"amount": .8,
					"name": "awesome",
				},

				want: "2awesome",
			},
			{

				name:       "Short-circuit OR",
				expression: "true || fail()",
				extension: Function("fail", func(arguments ...interface{}) (interface{}, error) {
					return nil, fmt.Errorf("Did not short-circuit")
				}),
				want: true,
			},
			{

				name:       "Short-circuit AND",
				expression: "false && fail()",
				extension: Function("fail", func(arguments ...interface{}) (interface{}, error) {
					return nil, fmt.Errorf("Did not short-circuit")
				}),
				want: false,
			},
			{

				name:       "Short-circuit ternary",
				expression: "true ? 1 : fail()",
				extension: Function("fail", func(arguments ...interface{}) (interface{}, error) {
					return nil, fmt.Errorf("Did not short-circuit")
				}),
				want: 1.0,
			},
			{

				name:       "Short-circuit coalesce",
				expression: `"foo" ?? fail()`,
				extension: Function("fail", func(arguments ...interface{}) (interface{}, error) {
					return nil, fmt.Errorf("Did not short-circuit")
				}),
				want: "foo",
			},
			{

				name:       "Simple parameter call",
				expression: "foo.String",
				parameter:  map[string]interface{}{"foo": foo},
				want:       foo.String,
			},
			{

				name:       "Simple parameter function call",
				expression: "foo.Func()",
				parameter:  map[string]interface{}{"foo": foo},
				want:       "funk",
			},
			{

				name:       "Simple parameter call from pointer",
				expression: "fooptr.String",
				parameter:  map[string]interface{}{"fooptr": &foo},
				want:       foo.String,
			},
			{

				name:       "Simple parameter function call from pointer",
				expression: "fooptr.Func()",
				parameter:  map[string]interface{}{"fooptr": &foo},
				want:       "funk",
			},
			{

				name:       "Simple parameter call",
				expression: `foo.String == "hi"`,
				parameter:  map[string]interface{}{"foo": foo},
				want:       false,
			},
			{

				name:       "Simple parameter call with modifier",
				expression: `foo.String + "hi"`,
				parameter:  map[string]interface{}{"foo": foo},
				want:       foo.String + "hi",
			},
			{

				name:       "Simple parameter function call, two-arg return",
				expression: `foo.Func2()`,
				parameter:  map[string]interface{}{"foo": foo},
				want:       "frink",
			},
			{

				name:       "Simple parameter function call, one arg",
				expression: `foo.FuncArgStr("boop")`,
				parameter:  map[string]interface{}{"foo": foo},
				want:       "boop",
			},
			{

				name:       "Simple parameter function call, one arg",
				expression: `foo.FuncArgStr("boop") + "hi"`,
				parameter:  map[string]interface{}{"foo": foo},
				want:       "boophi",
			},
			{

				name:       "Nested parameter function call",
				expression: `foo.Nested.Dunk("boop")`,
				parameter:  map[string]interface{}{"foo": foo},
				want:       "boopdunk",
			},
			{

				name:       "Nested parameter call",
				expression: "foo.Nested.Funk",
				parameter:  map[string]interface{}{"foo": foo},
				want:       "funkalicious",
			},
			{
				name:       "Nested map call",
				expression: `foo.Nested.Map["a"]`,
				parameter:  map[string]interface{}{"foo": foo},
				want:       1,
			},
			{
				name:       "Nested slice call",
				expression: `foo.Nested.Slice[1]`,
				parameter:  map[string]interface{}{"foo": foo},
				want:       2,
			},
			{

				name:       "Parameter call with + modifier",
				expression: "1 + foo.Int",
				parameter:  map[string]interface{}{"foo": foo},
				want:       102.0,
			},
			{

				name:       "Parameter string call with + modifier",
				expression: `"woop" + (foo.String)`,
				parameter:  map[string]interface{}{"foo": foo},
				want:       "woopstring!",
			},
			{

				name:       "Parameter call with && operator",
				expression: "true && foo.BoolFalse",
				parameter:  map[string]interface{}{"foo": foo},
				want:       false,
			},
			{
				name:       "Null coalesce nested parameter",
				expression: "foo.Nil ?? false",
				parameter:  map[string]interface{}{"foo": foo},
				want:       false,
			},
			{
				name:       "input functions",
				expression: "func1() + func2()",
				parameter: map[string]interface{}{
					"func1": func() float64 { return 2000 },
					"func2": func() float64 { return 2001 },
				},
				want: 4001.0,
			},
			{
				name:       "input functions",
				expression: "func1(date1) + func2(date2)",
				parameter: map[string]interface{}{
					"date1": func() interface{} {
						y2k, _ := time.Parse("2006", "2000")
						return y2k
					}(),
					"date2": func() interface{} {
						y2k1, _ := time.Parse("2006", "2001")
						return y2k1
					}(),
				},
				extension: NewLanguage(
					Function("func1", func(arguments ...interface{}) (interface{}, error) {
						return float64(arguments[0].(time.Time).Year()), nil
					}),
					Function("func2", func(arguments ...interface{}) (interface{}, error) {
						return float64(arguments[0].(time.Time).Year()), nil
					}),
				),
				want: 4001.0,
			},
			{
				name:       "complex64 number as parameter",
				expression: "complex64",
				parameter: map[string]interface{}{
					"complex64":  complex64(0),
					"complex128": complex128(0),
				},
				want: complex64(0),
			},
			{
				name:       "complex128 number as parameter",
				expression: "complex128",
				parameter: map[string]interface{}{
					"complex64":  complex64(0),
					"complex128": complex128(0),
				},
				want: complex128(0),
			},
			{
				name:       "coalesce with undefined",
				expression: "fooz ?? foo",
				parameter: map[string]interface{}{
					"foo": "bar",
				},
				want: "bar",
			},
			{
				name:       "map[interface{}]interface{}",
				expression: "foo",
				parameter: map[interface{}]interface{}{
					"foo": "bar",
				},
				want: "bar",
			},
			{
				name:       "method on pointer type",
				expression: "foo.PointerFunc()",
				parameter: map[string]interface{}{
					"foo": &dummyParameter{},
				},
				want: "point",
			},
			{
				name:       "custom selector",
				expression: "hello.world",
				parameter:  "!",
				extension: NewLanguage(Base(), VariableSelector(func(path Evaluables) Evaluable {
					return func(c context.Context, v interface{}) (interface{}, error) {
						keys, err := path.EvalStrings(c, v)
						if err != nil {
							return nil, err
						}
						return fmt.Sprintf("%s%s", strings.Join(keys, " "), v), nil
					}
				})),
				want: "hello world!",
			},
			{
				name:       "map[int]int",
				expression: `a[0] + a[2]`,
				parameter: map[string]interface{}{
					"a": map[int]int{0: 1, 2: 1},
				},
				want: 2.,
			},
			{
				name:       "map[int]string",
				expression: `a[0] * a[2]`,
				parameter: map[string]interface{}{
					"a": map[int]string{0: "1", 2: "1"},
				},
				want: 1.,
			},
		},
		t,
	)
}
