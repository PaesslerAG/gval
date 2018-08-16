package gval

// Courtesy of abrander
// ref: https://gist.github.com/abrander/fa05ae9b181b48ffe7afb12c961b6e90
import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var (
	hello  = "hello"
	empty  struct{}
	empty2 *string
	empty3 *int

	values = []interface{}{
		-1,
		0,
		12,
		13,
		"",
		"hello",
		&hello,
		nil,
		"nil",
		empty,
		empty2,
		true,
		false,
		time.Now(),
		rune('r'),
		int64(34),
		time.Duration(0),
		"true",
		"false",
		"\ntrue\n",
		"\nfalse\n",
		"12",
		"nil",
		"arg1",
		"arg2",
		int(12),
		int32(12),
		int64(12),
		complex(1.0, 1.0),
		[]byte{0, 0, 0},
		[]int{0, 0, 0},
		[]string{},
		"[]",
		"{}",
		"\"\"",
		"\"12\"",
		"\"hello\"",
		".*",
		"==",
		"!=",
		">",
		">=",
		"<",
		"<=",
		"=~",
		"!~",
		"in",
		"&&",
		"||",
		"^",
		"&",
		"|",
		">>",
		"<<",
		"+",
		"-",
		"*",
		"/",
		"%",
		"**",
		"-",
		"!",
		"~",
		"?",
		":",
		"??",
		"+",
		"-",
		"*",
		"/",
		"%",
		"**",
		"&",
		"|",
		"^",
		">>",
		"<<",
		",",
		"(",
		")",
		"[",
		"]",
		"\n",
		"\000",
	}

	panics = 0
)

const (
	SEED = 1487873697990155515
)

func BenchmarkRandom(bench *testing.B) {
	rand.Seed(SEED)
	for i := 0; i < bench.N; i++ {
		num := rand.Intn(3) + 2
		expression := ""

		for n := 0; n < num; n++ {
			expression += fmt.Sprintf(" %s", getRandom(values))
		}

		Evaluate(expression, nil)
	}
}

func getRandom(haystack []interface{}) interface{} {
	i := rand.Intn(len(haystack))
	return haystack[i]
}
