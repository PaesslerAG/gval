package gval

import (
	"strings"
	"testing"
	"text/scanner"
)

func TestMultilineStringLiterals(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       interface{}
		wantErr    string
	}{
		{
			name:       "Literal line feed in double quotes identical to escaped newline",
			expression: "\"hello\nworld\"",
			want:       "hello\nworld",
		},
		{
			name:       "Literal line feed comparison with escaped newline",
			expression: "\"hello\nworld\" == \"hello\\nworld\"",
			want:       true,
		},
		{
			name:       "Multiple literal line feeds",
			expression: "\"line 1\nline 2\nline 3\"",
			want:       "line 1\nline 2\nline 3",
		},
		{
			name:       "Literal carriage return and line feed",
			expression: "\"hello\r\nworld\"",
			want:       "hello\r\nworld",
		},
		{
			name:       "Multiline string with mixed escape sequences",
			expression: "\"first line \\\"quoted\\\"\t\nsecond line \\\\backslash\\\\\nthird line\"",
			want:       "first line \"quoted\"\t\nsecond line \\backslash\\\nthird line",
		},
		{
			name:       "Multiline string concatenation",
			expression: "\"hello\n\" + \"world\n\"",
			want:       "hello\nworld\n",
		},
		{
			name:       "Multiline string equality with raw backtick string",
			expression: "\"hello\nworld\" == `hello\nworld`",
			want:       true,
		},
		{
			name:       "Multiline string in array with 'in' operator",
			expression: "\"hello\nworld\" in [\"hello\nworld\", \"second\"]",
			want:       true,
		},
		{
			name:       "Multiline string in map object equality",
			expression: "{\"key\": \"val\nue\"} == {\"key\": \"val\nue\"}",
			want:       true,
		},
		{
			name:       "Multiline string in ternary expression",
			expression: "true ? \"yes\nline\" : \"no\"",
			want:       "yes\nline",
		},
		{
			name:       "Multiline string as function argument",
			expression: "date(\"2021-01-01\n\")",
			wantErr:    "could not parse",
		},
		{
			name:       "Unterminated multiline string returns parse error",
			expression: "\"hello\nworld",
			wantErr:    "could not parse string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.expression, nil)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("Evaluate(%q) expected error containing %q, got %v", tt.expression, tt.wantErr, got)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("Evaluate(%q) error %v does not contain %q", tt.expression, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Evaluate(%q) unexpected error: %v", tt.expression, err)
			}
			if got != tt.want {
				t.Fatalf("Evaluate(%q) = %v, want %v", tt.expression, got, tt.want)
			}
		})
	}
}

func TestMultilineStringPositionTracking(t *testing.T) {
	expression := "\"hello\nworld\"\n+ +"
	_, err := Evaluate(expression, nil)
	if err == nil {
		t.Fatal("expected syntax error on invalid operator sequence '+ +'")
	}
	errStr := err.Error()
	if !strings.Contains(errStr, ":3:") {
		t.Fatalf("expected error position on line 3, got: %s", errStr)
	}
}

func TestParserScanMultilineString(t *testing.T) {
	p := newParser("\"line 1\nline 2\"", Base())
	tok := p.Scan()
	if tok != scanner.String {
		t.Fatalf("p.Scan() = %v, want scanner.String", tok)
	}
	expectedToken := "\"line 1\nline 2\""
	if p.TokenText() != expectedToken {
		t.Fatalf("p.TokenText() = %q, want %q", p.TokenText(), expectedToken)
	}

	pos := p.scanner.Pos()
	if pos.Line != 2 {
		t.Fatalf("expected pos.Line == 2, got %d", pos.Line)
	}
}
