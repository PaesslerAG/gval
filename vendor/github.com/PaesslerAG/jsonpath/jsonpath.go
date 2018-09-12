package jsonpath

// jsonpath is an implementation of http://goessner.net/articles/JsonPath/
// If a JSONPath contains one of
// [key1, key2 ...], .., *, [min:max], [min:max:step], (? expression)
// all matchs are listed in an []interface{}
//
// The package comes with an extension of JSONPath to access the wildcard values of a match.
// If the JSONPath is used inside of a JSON object, you can use '#' or '#i' with natural number i
// to access all wildcards values or the ith wildcard

import (
	"context"

	"github.com/PaesslerAG/gval"
)

// New returns an selector for given jsonpath
func New(path string) (gval.Evaluable, error) {
	return lang.NewEvaluable(path)
}

//Get executes given jsonpath on given value
func Get(path string, value interface{}) (interface{}, error) {
	eval, err := lang.NewEvaluable(path)
	if err != nil {
		return nil, err
	}
	return eval(context.Background(), value)
}

var lang = gval.NewLanguage(
	gval.Base(),
	gval.PrefixExtension('$', single(getRootEvaluable).parse),
	gval.PrefixExtension('@', single(getCurrentEvaluable).parse),
)

//Language is the jsonpath Language
func Language() gval.Language {
	return lang
}

var wildcardExtension = gval.NewLanguage(
	lang,
	gval.PrefixExtension('{', parseJSONObject),
	gval.PrefixExtension('#', parseMatchReference),
)

//WildcardExtension is the jsonpath Language with access to the values, that matchs used wildcards
func WildcardExtension() gval.Language {
	return wildcardExtension
}
