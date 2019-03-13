package gval

import (
	"regexp/syntax"
	"testing"
)

func TestParsingFailure(t *testing.T) {
	testEvaluate(
		[]evaluationTest{
			{
				name:       "Invalid equality comparator",
				expression: "1 = 1",
				wantErr:    unexpected(`"="`, "operator"),
			},
			{
				name:       "Invalid equality comparator",
				expression: "1 === 1",
				wantErr:    unknownOp("==="),
			},
			{
				name:       "Too many characters for logical operator",
				expression: "true &&& false",
				wantErr:    unknownOp("&&&"),
			},
			{

				name:       "Too many characters for logical operator",
				expression: "true ||| false",
				wantErr:    unknownOp("|||"),
			},
			{

				name:       "Premature end to expression, via modifier",
				expression: "10 > 5 +",
				wantErr:    unexpected("EOF", "extensions"),
			},
			{
				name:       "Premature end to expression, via comparator",
				expression: "10 + 5 >",
				wantErr:    unexpected("EOF", "extensions"),
			},
			{
				name:       "Premature end to expression, via logical operator",
				expression: "10 > 5 &&",
				wantErr:    unexpected("EOF", "extensions"),
			},
			{

				name:       "Premature end to expression, via ternary operator",
				expression: "true ?",
				wantErr:    unexpected("EOF", "extensions"),
			},
			{
				name:       "Hanging REQ",
				expression: "`wat` =~",
				wantErr:    unexpected("EOF", "extensions"),
			},
			{

				name:       "Invalid operator change to REQ",
				expression: " / =~",
				wantErr:    unexpected(`"/"`, "extensions"),
			},
			{
				name:       "Invalid starting token, comparator",
				expression: "> 10",
				wantErr:    unexpected(`">"`, "extensions"),
			},
			{
				name:       "Invalid starting token, modifier",
				expression: "+ 5",
				wantErr:    unexpected(`"+"`, "extensions"),
			},
			{
				name:       "Invalid starting token, logical operator",
				expression: "&& 5 < 10",
				wantErr:    unexpected(`"&"`, "extensions"),
			},
			{
				name:       "Invalid NUMERIC transition",
				expression: "10 10",
				wantErr:    unexpected(`Int`, "operator"),
			},
			{
				name:       "Invalid STRING transition",
				expression: "`foo` `foo`",
				wantErr:    `String while scanning operator`, // can't use func unexpected because the token was changed from String to RawString in go 1.11
			},
			{
				name:       "Invalid operator transition",
				expression: "10 > < 10",
				wantErr:    unexpected(`"<"`, "extensions"),
			},
			{

				name:       "Starting with unbalanced parens",
				expression: " ) ( arg2",
				wantErr:    unexpected(`")"`, "extensions"),
			},
			{

				name:       "Unclosed bracket",
				expression: "[foo bar",
				wantErr:    unexpected(`EOF`, "extensions"),
			},
			{

				name:       "Unclosed quote",
				expression: "foo == `responseTime",
				wantErr:    "could not parse string",
			},
			{

				name:       "Constant regex pattern fail to compile",
				expression: "foo =~ `[abc`",
				wantErr:    string(syntax.ErrMissingBracket),
			},
			{

				name:       "Constant unmatch regex pattern fail to compile",
				expression: "foo !~ `[abc`",
				wantErr:    string(syntax.ErrMissingBracket),
			},
			{

				name:       "Unbalanced parentheses",
				expression: "10 > (1 + 50",
				wantErr:    unexpected(`EOF`, "parentheses"),
			},
			{

				name:       "Multiple radix",
				expression: "127.0.0.1",
				wantErr:    unexpected(`Float`, "operator"),
			},
			{

				name:       "Hanging accessor",
				expression: "foo.Bar.",
				wantErr:    unexpected(`EOF`, "field"),
			},
			{
				name:       "Incomplete Hex",
				expression: "0x",
				wantErr:    `strconv.ParseFloat: parsing "0x": invalid syntax`,
			},
			{
				name:       "Invalid Hex literal",
				expression: "0x > 0",
				wantErr:    `strconv.ParseFloat: parsing "0x": invalid syntax`,
			},
			{
				name:       "Hex float (Unsupported)",
				expression: "0x1.1",
				wantErr:    `strconv.ParseFloat: parsing "0x1": invalid syntax`,
			},
			{
				name:       "Hex invalid letter",
				expression: "0x12g1",
				wantErr:    `strconv.ParseFloat: parsing "0x12": invalid syntax`,
			},
			{
				name:       "Error after camouflage",
				expression: "0 + ,",
				wantErr:    `unexpected "," while scanning extensions`,
			},
		},
		t,
	)
}

func unknownOp(op string) string {
	return "unknown operator " + op
}

func unexpected(token, unit string) string {
	return "unexpected " + token + " while scanning " + unit
}
