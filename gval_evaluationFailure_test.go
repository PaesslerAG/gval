package gval

/*
	Tests to make sure evaluation fails in the expected ways.
*/
import (
	"errors"
	"fmt"
	"testing"
)

func TestModifierTyping(test *testing.T) {
	var (
		invalidOperator      = "invalid operation"
		unknownParameter     = "unknown parameter"
		invalidRegex         = "error parsing regex"
		tooFewArguments      = "reflect: Call with too few input arguments"
		tooManyArguments     = "reflect: Call with too many input arguments"
		mismatchedParameters = "reflect: Call using"
		custom               = "test error"
	)
	evaluationTests := []evaluationTest{
		//ModifierTyping
		{
			name:       "PLUS literal number to literal bool",
			expression: "1 + true",
			want:       "1true", // + on string is defined
		},
		{
			name:       "PLUS number to bool",
			expression: "number + bool",
			want:       "1true", // + on string is defined
		},
		{
			name:       "MINUS number to bool",
			expression: "number - bool",
			wantErr:    invalidOperator,
		},
		{
			name:       "MINUS number to bool",
			expression: "number - bool",
			wantErr:    invalidOperator,
		},
		{
			name:       "MULTIPLY number to bool",
			expression: "number * bool",
			wantErr:    invalidOperator,
		},
		{
			name:       "DIVIDE number to bool",
			expression: "number / bool",
			wantErr:    invalidOperator,
		},
		{
			name:       "EXPONENT number to bool",
			expression: "number ** bool",
			wantErr:    invalidOperator,
		},
		{
			name:       "MODULUS number to bool",
			expression: "number % bool",
			wantErr:    invalidOperator,
		},
		{
			name:       "XOR number to bool",
			expression: "number % bool",
			wantErr:    invalidOperator,
		},
		{
			name:       "BITWISE_OR number to bool",
			expression: "number | bool",
			wantErr:    invalidOperator,
		},
		{
			name:       "BITWISE_AND number to bool",
			expression: "number & bool",
			wantErr:    invalidOperator,
		},
		{
			name:       "BITWISE_XOR number to bool",
			expression: "number ^ bool",
			wantErr:    invalidOperator,
		},
		{
			name:       "BITWISE_LSHIFT number to bool",
			expression: "number << bool",
			wantErr:    invalidOperator,
		},
		{
			name:       "BITWISE_RSHIFT number to bool",
			expression: "number >> bool",
			wantErr:    invalidOperator,
		},
		//LogicalOperatorTyping
		{
			name:       "AND number to number",
			expression: "number && number",
			want:       true, // number != 0 is true
		},
		{

			name:       "OR number to number",
			expression: "number || number",
			want:       true, // number != 0 is true
		},
		{
			name:       "AND string to string",
			expression: "string && string",
			wantErr:    invalidOperator,
		},
		{
			name:       "OR string to string",
			expression: "string || string",
			wantErr:    invalidOperator,
		},
		{
			name:       "AND number to string",
			expression: "number && string",
			wantErr:    invalidOperator,
		},
		{
			name:       "OR number to string",
			expression: "number || string",
			wantErr:    invalidOperator,
		},
		{
			name:       "AND bool to string",
			expression: "bool && string",
			wantErr:    invalidOperator,
		},
		{
			name:       "OR string to bool",
			expression: "string || bool",
			wantErr:    invalidOperator,
		},
		//ComparatorTyping
		{
			name:       "GT literal bool to literal bool",
			expression: "true > true",
			want:       false, //lexical order on "true"
		},
		{
			name:       "GT bool to bool",
			expression: "bool > bool",
			want:       false, //lexical order on "true"
		},
		{
			name:       "GTE bool to bool",
			expression: "bool >= bool",
			want:       true, //lexical order on "true"
		},
		{
			name:       "LT bool to bool",
			expression: "bool < bool",
			want:       false, //lexical order on "true"
		},
		{
			name:       "LTE bool to bool",
			expression: "bool <= bool",
			want:       true, //lexical order on "true"
		},
		{
			name:       "GT number to string",
			expression: "number > string",
			want:       false, //lexical order "1" < "foo"
		},
		{

			name:       "GTE number to string",
			expression: "number >= string",
			want:       false, //lexical order "1" < "foo"
		},
		{
			name:       "LT number to string",
			expression: "number < string",
			want:       true, //lexical order "1" < "foo"
		},
		{
			name:       "REQ number to string",
			expression: "number =~ string",
			want:       false,
		},
		{
			name:       "REQ number to bool",
			expression: "number =~ bool",
			want:       false,
		},
		{
			name:       "REQ bool to number",
			expression: "bool =~ number",
			want:       false,
		},
		{
			name:       "REQ bool to string",
			expression: "bool =~ string",
			want:       false,
		},
		{
			name:       "NREQ number to string",
			expression: "number !~ string",
			want:       true,
		},
		{
			name:       "NREQ number to bool",
			expression: "number !~ bool",
			want:       true,
		},
		{
			name:       "NREQ bool to number",
			expression: "bool !~ number",
			want:       true,
		},
		{

			name:       "NREQ bool to string",
			expression: "bool !~ string",
			want:       true,
		},
		{
			name:       "IN non-array numeric",
			expression: "1 in 2",
			wantErr:    "expected type []interface{} for in operator but got float64",
		},
		{
			name:       "IN non-array string",
			expression: `1 in "foo"`,
			wantErr:    "expected type []interface{} for in operator but got string",
		},
		{

			name:       "IN non-array boolean",
			expression: "1 in true",
			wantErr:    "expected type []interface{} for in operator but got bool",
		},
		//TernaryTyping
		{
			name:       "Ternary with number",
			expression: "10 ? true",
			want:       true, // 10 != nil && 10 != false
		},
		{
			name:       "Ternary with string",
			expression: `"foo" ? true`,
			want:       true, // "foo" != nil && "foo" != false
		},
		//RegexParameterCompilation
		{
			name:       "Regex equality runtime parsing",
			expression: `"foo" =~ foo`,
			parameter: map[string]interface{}{
				"foo": "[foo",
			},
			wantErr: invalidRegex,
		},
		{
			name:       "Regex inequality runtime parsing",
			expression: `"foo" !~ foo`,
			parameter: map[string]interface{}{
				"foo": "[foo",
			},
			wantErr: invalidRegex,
		},
		{
			name:       "Regex equality runtime right side evaluation",
			expression: `"foo" =~ error()`,
			wantErr:    custom,
		},
		{
			name:       "Regex inequality runtime right side evaluation",
			expression: `"foo" !~ error()`,
			wantErr:    custom,
		},
		{
			name:       "Regex equality runtime left side evaluation",
			expression: `error() =~ "."`,
			wantErr:    custom,
		},
		{
			name:       "Regex inequality runtime left side evaluation",
			expression: `error() !~ "."`,
			wantErr:    custom,
		},
		//FuncExecution
		{
			name:       "Func error bubbling",
			expression: "error()",
			extension: Function("error", func(arguments ...interface{}) (interface{}, error) {
				return nil, errors.New("Huge problems")
			}),
			wantErr: "Huge problems",
		},
		//InvalidParameterCalls
		{
			name:       "Missing parameter field reference",
			expression: "foo.NotExists",
			parameter:  fooFailureParameters,
			wantErr:    unknownParameter,
		},
		{
			name:       "Parameter method call on missing function",
			expression: "foo.NotExist()",
			parameter:  fooFailureParameters,
			wantErr:    unknownParameter,
		},
		{
			name:       "Nested missing parameter field reference",
			expression: "foo.Nested.NotExists",
			parameter:  fooFailureParameters,
			wantErr:    unknownParameter,
		},
		{
			name:       "Parameter method call returns error",
			expression: "foo.AlwaysFail()",
			parameter:  fooFailureParameters,
			wantErr:    "function should always fail",
		},
		{
			name:       "Too few arguments to parameter call",
			expression: "foo.FuncArgStr()",
			parameter:  fooFailureParameters,
			wantErr:    tooFewArguments,
		},
		{
			name:       "Too many arguments to parameter call",
			expression: `foo.FuncArgStr("foo", "bar", 15)`,
			parameter:  fooFailureParameters,
			wantErr:    tooManyArguments,
		},
		{
			name:       "Mismatched parameters",
			expression: "foo.FuncArgStr(5)",
			parameter:  fooFailureParameters,
			wantErr:    mismatchedParameters,
		},
		{
			name:       "Negative Array Index",
			expression: "foo[-1]",
			parameter: map[string]interface{}{
				"foo": []int{1, 2, 3},
			},
			wantErr: unknownParameter,
		},
		{
			name:       "Nested slice call index out of bound",
			expression: `foo.Nested.Slice[10]`,
			parameter:  map[string]interface{}{"foo": foo},
			wantErr:    unknownParameter,
		},
		{
			name:       "Nested map call missing key",
			expression: `foo.Nested.Map["d"]`,
			parameter:  map[string]interface{}{"foo": foo},
			wantErr:    unknownParameter,
		},
		{
			name:       "invalid selector",
			expression: "hello[world()]",
			extension: NewLanguage(Base(), Function("world", func() (int, error) {
				return 0, fmt.Errorf("test error")
			})),
			wantErr: "test error",
		},
		{
			name:       "eval `nil > 1` returns true #23",
			expression: `nil > 1`,
			wantErr:    "invalid operation (<nil>) > (float64)",
		},
	}

	for i := range evaluationTests {
		if evaluationTests[i].parameter == nil {
			evaluationTests[i].parameter = map[string]interface{}{
				"number": 1,
				"string": "foo",
				"bool":   true,
				"error": func() (int, error) {
					return 0, fmt.Errorf("test error")
				},
			}
		}
	}

	testEvaluate(evaluationTests, test)
}
