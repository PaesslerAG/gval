package gval_test

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/gval"
	"github.com/PaesslerAG/jsonpath"
)

func Example() {

	vars := map[string]interface{}{"name": "World"}

	value, err := gval.Evaluate(`"Hello " + name + "!"`, vars)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(value)

	// Output:
	// Hello World!
}

func ExampleEvaluate() {

	value, err := gval.Evaluate("foo > 0", map[string]interface{}{
		"foo": -1.,
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(value)

	// Output:
	// false
}

func ExampleEvaluate_nestedParameter() {

	value, err := gval.Evaluate("foo.bar > 0", map[string]interface{}{
		"foo": map[string]interface{}{"bar": -1.},
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(value)

	// Output:
	// false
}

func ExampleEvaluate_array() {

	value, err := gval.Evaluate("foo[0]", map[string]interface{}{
		"foo": []interface{}{-1.},
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(value)

	// Output:
	// -1
}

func ExampleEvaluate_complexAccessor() {

	value, err := gval.Evaluate(`foo["b" + "a" + "r"]`, map[string]interface{}{
		"foo": map[string]interface{}{"bar": -1.},
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(value)

	// Output:
	// -1
}

func ExampleEvaluate_arithmetic() {

	value, err := gval.Evaluate("(requests_made * requests_succeeded / 100) >= 90",
		map[string]interface{}{
			"requests_made":      100,
			"requests_succeeded": 80,
		})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(value)

	// Output:
	// false
}

func ExampleEvaluate_string() {

	value, err := gval.Evaluate(`http_response_body == "service is ok"`,
		map[string]interface{}{
			"http_response_body": "service is ok",
		})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(value)

	// Output:
	// true
}

func ExampleEvaluate_float64() {

	value, err := gval.Evaluate("(mem_used / total_mem) * 100",
		map[string]interface{}{
			"total_mem": 1024,
			"mem_used":  512,
		})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(value)

	// Output:
	// 50
}

func ExampleEvaluate_dateComparison() {

	value, err := gval.Evaluate("date(`2014-01-02`) > date(`2014-01-01 23:59:59`)",
		nil,
		// define Date comparison because it is not part expression language gval
		gval.InfixOperator(">", func(a, b interface{}) (interface{}, error) {
			date1, ok1 := a.(time.Time)
			date2, ok2 := b.(time.Time)

			if ok1 && ok2 {
				return date1.After(date2), nil
			}
			return nil, fmt.Errorf("unexpected operands types (%T) > (%T)", a, b)
		}),
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(value)

	// Output:
	// true
}

func ExampleEvaluable() {
	eval, err := gval.Full(gval.Constant("maximum_time", 52)).
		NewEvaluable("response_time <= maximum_time")
	if err != nil {
		fmt.Println(err)
	}

	for i := 50; i < 55; i++ {
		value, err := eval(context.Background(), map[string]interface{}{
			"response_time": i,
		})
		if err != nil {
			fmt.Println(err)

		}

		fmt.Println(value)
	}

	// Output:
	// true
	// true
	// true
	// false
	// false
}

func ExampleEvaluate_strlen() {

	value, err := gval.Evaluate(`strlen("someReallyLongInputString") <= 16`,
		nil,
		gval.Function("strlen", func(args ...interface{}) (interface{}, error) {
			length := len(args[0].(string))
			return (float64)(length), nil
		}))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(value)

	// Output:
	// false
}

func ExampleEvaluate_encoding() {

	value, err := gval.Evaluate(`(7 < "47" == true ? "hello world!\n\u263a" : "good bye\n")`+" + ` more text`",
		nil,
		gval.Function("strlen", func(args ...interface{}) (interface{}, error) {
			length := len(args[0].(string))
			return (float64)(length), nil
		}))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(value)

	// Output:
	// hello world!
	// ☺ more text
}

type exampleType struct {
	Hello string
}

func (e exampleType) World() string {
	return "world"
}

func ExampleEvaluate_accessor() {

	value, err := gval.Evaluate(`foo.Hello + foo.World()`,
		map[string]interface{}{
			"foo": exampleType{Hello: "hello "},
		})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(value)

	// Output:
	// hello world
}

func ExampleEvaluate_flatAccessor() {

	value, err := gval.Evaluate(`Hello + World()`,
		exampleType{Hello: "hello "},
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(value)

	// Output:
	// hello world
}

func ExampleEvaluate_nestedAccessor() {

	value, err := gval.Evaluate(`foo.Bar.Hello + foo.Bar.World()`,
		map[string]interface{}{
			"foo": struct{ Bar exampleType }{
				Bar: exampleType{Hello: "hello "},
			},
		})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(value)

	// Output:
	// hello world
}

func ExampleVariableSelector() {
	value, err := gval.Evaluate(`hello.world`,
		"!",
		gval.VariableSelector(func(path gval.Evaluables) gval.Evaluable {
			return func(c context.Context, v interface{}) (interface{}, error) {
				keys, err := path.EvalStrings(c, v)
				if err != nil {
					return nil, err
				}
				return fmt.Sprintf("%s%s", strings.Join(keys, " "), v), nil
			}
		}),
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(value)

	// Output:
	// hello world!
}

func ExampleEvaluable_EvalInt() {
	eval, err := gval.Full().NewEvaluable("1 + x")
	if err != nil {
		fmt.Println(err)
		return
	}

	value, err := eval.EvalInt(context.Background(), map[string]interface{}{"x": 5})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(value)

	// Output:
	// 6
}

func ExampleEvaluable_EvalBool() {
	eval, err := gval.Full().NewEvaluable("1 == x")
	if err != nil {
		fmt.Println(err)
		return
	}

	value, err := eval.EvalBool(context.Background(), map[string]interface{}{"x": 1})
	if err != nil {
		fmt.Println(err)
	}

	if value {
		fmt.Print("yeah")
	}

	// Output:
	// yeah
}

func ExampleEvaluate_jsonpath() {

	value, err := gval.Evaluate(`$["response-time"]`,
		map[string]interface{}{
			"response-time": 100,
		},
		jsonpath.Language(),
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(value)

	// Output:
	// 100
}

func ExampleLanguage() {
	lang := gval.NewLanguage(gval.JSON(), gval.Arithmetic(),
		//pipe operator
		gval.PostfixOperator("|", func(c context.Context, p *gval.Parser, pre gval.Evaluable) (gval.Evaluable, error) {
			post, err := p.ParseExpression(c)
			if err != nil {
				return nil, err
			}
			return func(c context.Context, v interface{}) (interface{}, error) {
				v, err := pre(c, v)
				if err != nil {
					return nil, err
				}
				return post(c, v)
			}, nil
		}))

	eval, err := lang.NewEvaluable(`{"foobar": 50} | foobar + 100`)
	if err != nil {
		fmt.Println(err)
	}

	value, err := eval(context.Background(), nil)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(value)

	// Output:
	// 150
}
