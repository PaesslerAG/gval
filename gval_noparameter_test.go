package gval

import (
	"context"
	"fmt"
	"testing"
)

func TestNoParameter(t *testing.T) {
	testEvaluate(
		[]evaluationTest{
			{
				name:       "Number",
				expression: "100",
				want:       100.0,
			},
			{
				name:       "Single PLUS",
				expression: "51 + 49",
				want:       100.0,
			},
			{
				name:       "Single MINUS",
				expression: "100 - 51",
				want:       49.0,
			},
			{
				name:       "Single BITWISE AND",
				expression: "100 & 50",
				want:       32.0,
			},
			{
				name:       "Single BITWISE OR",
				expression: "100 | 50",
				want:       118.0,
			},
			{
				name:       "Single BITWISE XOR",
				expression: "100 ^ 50",
				want:       86.0,
			},
			{
				name:       "Single shift left",
				expression: "2 << 1",
				want:       4.0,
			},
			{
				name:       "Single shift right",
				expression: "2 >> 1",
				want:       1.0,
			},
			{
				name:       "Single BITWISE NOT",
				expression: "~10",
				want:       -11.0,
			},
			{

				name:       "Single MULTIPLY",
				expression: "5 * 20",
				want:       100.0,
			},
			{

				name:       "Single DIVIDE",
				expression: "100 / 20",
				want:       5.0,
			},
			{

				name:       "Single even MODULUS",
				expression: "100 % 2",
				want:       0.0,
			},
			{
				name:       "Single odd MODULUS",
				expression: "101 % 2",
				want:       1.0,
			},
			{

				name:       "Single EXPONENT",
				expression: "10 ** 2",
				want:       100.0,
			},
			{

				name:       "Compound PLUS",
				expression: "20 + 30 + 50",
				want:       100.0,
			},
			{

				name:       "Compound BITWISE AND",
				expression: "20 & 30 & 50",
				want:       16.0,
			},
			{
				name:       "Mutiple operators",
				expression: "20 * 5 - 49",
				want:       51.0,
			},
			{
				name:       "Parenthesis usage",
				expression: "100 - (5 * 10)",
				want:       50.0,
			},
			{

				name:       "Nested parentheses",
				expression: "50 + (5 * (15 - 5))",
				want:       100.0,
			},
			{

				name:       "Nested parentheses with bitwise",
				expression: "100 ^ (23 * (2 | 5))",
				want:       197.0,
			},
			{
				name:       "Logical OR operation of two clauses",
				expression: "(1 == 1) || (true == true)",
				want:       true,
			},
			{
				name:       "Logical AND operation of two clauses",
				expression: "(1 == 1) && (true == true)",
				want:       true,
			},
			{

				name:       "Implicit boolean",
				expression: "2 > 1",
				want:       true,
			},
			{

				name:       "Compound boolean",
				expression: "5 < 10 && 1 < 5",
				want:       true,
			},
			{
				name:       "Evaluated true && false operation (for issue #8)",
				expression: "1 > 10 && 11 > 10",
				want:       false,
			},
			{

				name:       "Evaluated true && false operation (for issue #8)",
				expression: "true == true && false == true",
				want:       false,
			},
			{

				name:       "Parenthesis boolean",
				expression: "10 < 50 && (1 != 2 && 1 > 0)",
				want:       true,
			},
			{
				name:       "Comparison of string constants",
				expression: `"foo" == "foo"`,
				want:       true,
			},
			{

				name:       "NEQ comparison of string constants",
				expression: `"foo" != "bar"`,
				want:       true,
			},
			{

				name:       "REQ comparison of string constants",
				expression: `"foobar" =~ "oba"`,
				want:       true,
			},
			{

				name:       "NREQ comparison of string constants",
				expression: `"foo" !~ "bar"`,
				want:       true,
			},
			{

				name:       "Multiplicative/additive order",
				expression: "5 + 10 * 2",
				want:       25.0,
			},
			{
				name:       "Multiple constant multiplications",
				expression: "10 * 10 * 10",
				want:       1000.0,
			},
			{

				name:       "Multiple adds/multiplications",
				expression: "10 * 10 * 10 + 1 * 10 * 10",
				want:       1100.0,
			},
			{

				name:       "Modulus operatorPrecedence",
				expression: "1 + 101 % 2 * 5",
				want:       6.0,
			},
			{
				name:       "Exponent operatorPrecedence",
				expression: "1 + 5 ** 3 % 2 * 5",
				want:       6.0,
			},
			{

				name:       "Bit shift operatorPrecedence",
				expression: "50 << 1 & 90",
				want:       64.0,
			},
			{

				name:       "Bit shift operatorPrecedence",
				expression: "90 & 50 << 1",
				want:       64.0,
			},
			{

				name:       "Bit shift operatorPrecedence amongst non-bitwise",
				expression: "90 + 50 << 1 * 5",
				want:       4480.0,
			},
			{
				name:       "Order of non-commutative same-operatorPrecedence operators (additive)",
				expression: "1 - 2 - 4 - 8",
				want:       -13.0,
			},
			{
				name:       "Order of non-commutative same-operatorPrecedence operators (multiplicative)",
				expression: "1 * 4 / 2 * 8",
				want:       16.0,
			},
			{
				name:       "Null coalesce operatorPrecedence",
				expression: "true ?? true ? 100 + 200 : 400",
				want:       300.0,
			},
			{
				name:       "Identical date equivalence",
				expression: `"2014-01-02 14:12:22" == "2014-01-02 14:12:22"`,
				want:       true,
			},
			{
				name:       "Positive date GT",
				expression: `"2014-01-02 14:12:22" > "2014-01-02 12:12:22"`,
				want:       true,
			},
			{
				name:       "Negative date GT",
				expression: `"2014-01-02 14:12:22" > "2014-01-02 16:12:22"`,
				want:       false,
			},
			{
				name:       "Positive date GTE",
				expression: `"2014-01-02 14:12:22" >= "2014-01-02 12:12:22"`,
				want:       true,
			},
			{
				name:       "Negative date GTE",
				expression: `"2014-01-02 14:12:22" >= "2014-01-02 16:12:22"`,
				want:       false,
			},
			{
				name:       "Positive date LT",
				expression: `"2014-01-02 14:12:22" < "2014-01-02 16:12:22"`,
				want:       true,
			},
			{

				name:       "Negative date LT",
				expression: `"2014-01-02 14:12:22" < "2014-01-02 11:12:22"`,
				want:       false,
			},
			{

				name:       "Positive date LTE",
				expression: `"2014-01-02 09:12:22" <= "2014-01-02 12:12:22"`,
				want:       true,
			},
			{
				name:       "Negative date LTE",
				expression: `"2014-01-02 14:12:22" <= "2014-01-02 11:12:22"`,
				want:       false,
			},
			{

				name:       "Sign prefix comparison",
				expression: "-1 < 0",
				want:       true,
			},
			{

				name:       "Lexicographic LT",
				expression: `"ab" < "abc"`,
				want:       true,
			},
			{
				name:       "Lexicographic LTE",
				expression: `"ab" <= "abc"`,
				want:       true,
			},
			{

				name:       "Lexicographic GT",
				expression: `"aba" > "abc"`,
				want:       false,
			},
			{

				name:       "Lexicographic GTE",
				expression: `"aba" >= "abc"`,
				want:       false,
			},
			{

				name:       "Boolean sign prefix comparison",
				expression: "!true == false",
				want:       true,
			},
			{
				name:       "Inversion of clause",
				expression: "!(10 < 0)",
				want:       true,
			},
			{

				name:       "Negation after modifier",
				expression: "10 * -10",
				want:       -100.0,
			},
			{

				name:       "Ternary with single boolean",
				expression: "true ? 10",
				want:       10.0,
			},
			{

				name:       "Ternary nil with single boolean",
				expression: "false ? 10",
				want:       nil,
			},
			{
				name:       "Ternary with comparator boolean",
				expression: "10 > 5 ? 35.50",
				want:       35.50,
			},
			{

				name:       "Ternary nil with comparator boolean",
				expression: "1 > 5 ? 35.50",
				want:       nil,
			},
			{

				name:       "Ternary with parentheses",
				expression: "(5 * (15 - 5)) > 5 ? 35.50",
				want:       35.50,
			},
			{

				name:       "Ternary operatorPrecedence",
				expression: "true ? 35.50 > 10",
				want:       true,
			},
			{
				name:       "Ternary-else",
				expression: "false ? 35.50 : 50",
				want:       50.0,
			},
			{

				name:       "Ternary-else inside clause",
				expression: "(false ? 5 : 35.50) > 10",
				want:       true,
			},
			{

				name:       "Ternary-else (true-case) inside clause",
				expression: "(true ? 1 : 5) < 10",
				want:       true,
			},
			{

				name:       "Ternary-else before comparator (negative case)",
				expression: "true ? 1 : 5 > 10",
				want:       1.0,
			},
			{
				name:       "Nested ternaries (#32)",
				expression: "(2 == 2) ? 1 : (true ? 2 : 3)",
				want:       1.0,
			},
			{

				name:       "Nested ternaries, right case (#32)",
				expression: "false ? 1 : (true ? 2 : 3)",
				want:       2.0,
			},
			{

				name:       "Doubly-nested ternaries (#32)",
				expression: "true ? (false ? 1 : (false ? 2 : 3)) : (false ? 4 : 5)",
				want:       3.0,
			},
			{

				name:       "String to string concat",
				expression: `"foo" + "bar" == "foobar"`,
				want:       true,
			},
			{
				name:       "String to float64 concat",
				expression: `"foo" + 123 == "foo123"`,
				want:       true,
			},
			{

				name:       "Float64 to string concat",
				expression: `123 + "bar" == "123bar"`,
				want:       true,
			},
			{

				name:       "String to date concat",
				expression: `"foo" + "02/05/1970" == "foobar"`,
				want:       false,
			},
			{

				name:       "String to bool concat",
				expression: `"foo" + true == "footrue"`,
				want:       true,
			},
			{
				name:       "Bool to string concat",
				expression: `true + "bar" == "truebar"`,
				want:       true,
			},
			{

				name:       "Null coalesce left",
				expression: "1 ?? 2",
				want:       1.0,
			},
			{

				name:       "Array membership literals",
				expression: "1 in [1, 2, 3]",
				want:       true,
			},
			{

				name:       "Array membership literal with inversion",
				expression: "!(1 in [1, 2, 3])",
				want:       false,
			},
			{
				name:       "Logical operator reordering (#30)",
				expression: "(true && true) || (true && false)",
				want:       true,
			},
			{

				name:       "Logical operator reordering without parens (#30)",
				expression: "true && true || true && false",
				want:       true,
			},
			{

				name:       "Logical operator reordering with multiple OR (#30)",
				expression: "false || true && true || false",
				want:       true,
			},
			{
				name:       "Left-side multiple consecutive (should be reordered) operators",
				expression: "(10 * 10 * 10) > 10",
				want:       true,
			},
			{

				name:       "Three-part non-paren logical op reordering (#44)",
				expression: "false && true || true",
				want:       true,
			},
			{

				name:       "Three-part non-paren logical op reordering (#44), second one",
				expression: "true || false && true",
				want:       true,
			},
			{
				name:       "Logical operator reordering without parens (#45)",
				expression: "true && true || false && false",
				want:       true,
			},
			{
				name:       "Single function",
				expression: "foo()",
				extension: Function("foo", func(arguments ...interface{}) (interface{}, error) {
					return true, nil
				}),

				want: true,
			},
			{
				name:       "Func with argument",
				expression: "passthrough(1)",
				extension: Function("passthrough", func(arguments ...interface{}) (interface{}, error) {
					return arguments[0], nil
				}),
				want: 1.0,
			},
			{
				name:       "Func with arguments",
				expression: "passthrough(1, 2)",
				extension: Function("passthrough", func(arguments ...interface{}) (interface{}, error) {
					return arguments[0].(float64) + arguments[1].(float64), nil
				}),
				want: 3.0,
			},
			{
				name:       "Nested function with operatorPrecedence",
				expression: "sum(1, sum(2, 3), 2 + 2, true ? 4 : 5)",
				extension: Function("sum", func(arguments ...interface{}) (interface{}, error) {
					sum := 0.0
					for _, v := range arguments {
						sum += v.(float64)
					}
					return sum, nil
				}),
				want: 14.0,
			},
			{
				name:       "Empty function and modifier, compared",
				expression: "numeric()-1 > 0",
				extension: Function("numeric", func(arguments ...interface{}) (interface{}, error) {
					return 2.0, nil
				}),
				want: true,
			},
			{
				name:       "Empty function comparator",
				expression: "numeric() > 0",
				extension: Function("numeric", func(arguments ...interface{}) (interface{}, error) {
					return 2.0, nil
				}),
				want: true,
			},
			{

				name:       "Empty function logical operator",
				expression: "success() && !false",
				extension: Function("success", func(arguments ...interface{}) (interface{}, error) {
					return true, nil
				}),
				want: true,
			},
			{
				name:       "Empty function ternary",
				expression: "nope() ? 1 : 2.0",
				extension: Function("nope", func(arguments ...interface{}) (interface{}, error) {
					return false, nil
				}),
				want: 2.0,
			},
			{

				name:       "Empty function null coalesce",
				expression: "null() ?? 2",
				extension: Function("null", func(arguments ...interface{}) (interface{}, error) {
					return nil, nil
				}),
				want: 2.0,
			},
			{
				name:       "Empty function with prefix",
				expression: "-ten()",
				extension: Function("ten", func(arguments ...interface{}) (interface{}, error) {
					return 10.0, nil
				}),
				want: -10.0,
			},
			{
				name:       "Empty function as part of chain",
				expression: "10 - numeric() - 2",
				extension: Function("numeric", func(arguments ...interface{}) (interface{}, error) {
					return 5.0, nil
				}),
				want: 3.0,
			},
			{
				name:       "Empty function near separator",
				expression: "10 in [1, 2, 3, ten(), 8]",
				extension: Function("ten", func(arguments ...interface{}) (interface{}, error) {
					return 10.0, nil
				}),
				want: true,
			},
			{
				name:       "Enclosed empty function with modifier and comparator (#28)",
				expression: "(ten() - 1) > 3",
				extension: Function("ten", func(arguments ...interface{}) (interface{}, error) {
					return 10.0, nil
				}),
				want: true,
			},
			{
				name:       "Array",
				expression: `[(ten() - 1) > 3, (ten() - 1),"hey"]`,
				extension: Function("ten", func(arguments ...interface{}) (interface{}, error) {
					return 10.0, nil
				}),
				want: []interface{}{true, 9., "hey"},
			},
			{
				name:       "Object",
				expression: `{1: (ten() - 1) > 3, 7 + ".X" : (ten() - 1),"hello" : "hey"}`,
				extension: Function("ten", func(arguments ...interface{}) (interface{}, error) {
					return 10.0, nil
				}),
				want: map[string]interface{}{"1": true, "7.X": 9., "hello": "hey"},
			},
			{
				name:       "Object negativ value",
				expression: `{1: -1,"hello" : "hey"}`,
				want:       map[string]interface{}{"1": -1., "hello": "hey"},
			},
			{
				name:       "Empty Array",
				expression: `[]`,
				want:       []interface{}{},
			},
			{
				name:       "Empty Object",
				expression: `{}`,
				want:       map[string]interface{}{},
			},
			{
				name:       "Variadic",
				expression: `sum(1,2,3,4)`,
				extension: Function("sum", func(arguments ...float64) (interface{}, error) {
					sum := 0.
					for _, a := range arguments {
						sum += a
					}
					return sum, nil
				}),
				want: 10.0,
			},
			{
				name:       "Ident Operator",
				expression: `1 plus 1`,
				extension: InfixNumberOperator("plus", func(a, b float64) (interface{}, error) {
					return a + b, nil
				}),
				want: 2.0,
			},
			{
				name:       "Postfix Operator",
				expression: `4§`,
				extension: PostfixOperator("§", func(_ context.Context, _ *Parser, eval Evaluable) (Evaluable, error) {
					return func(ctx context.Context, parameter interface{}) (interface{}, error) {
						i, err := eval.EvalInt(ctx, parameter)
						if err != nil {
							return nil, err
						}
						return fmt.Sprintf("§%d", i), nil
					}, nil
				}),
				want: "§4",
			},
		},
		t,
	)
}
