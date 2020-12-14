package gval

import (
	"testing"
	"text/scanner"
	"unicode"
)

func TestParser_Scan(t *testing.T) {
	tests := []struct {
		name  string
		input string
		Language
		do        func(p *Parser)
		wantScan  rune
		wantToken string
		wantPanic bool
	}{
		{
			name:  "camouflage",
			input: "$abc",
			do: func(p *Parser) {
				p.Scan()
				p.Camouflage("test")
			},
			wantScan:  '$',
			wantToken: "$",
		},
		{
			name:  "camouflage with next",
			input: "$abc",
			do: func(p *Parser) {
				p.Scan()
				p.Camouflage("test")
				p.Next()
			},
			wantPanic: true,
		},
		{
			name:  "camouflage scan camouflage",
			input: "$abc",
			do: func(p *Parser) {
				p.Scan()
				p.Camouflage("test")
				p.Scan()
				p.Camouflage("test2")
			},
			wantScan:  '$',
			wantToken: "$",
		},
		{
			name:  "camouflage with peek",
			input: "$abc",
			do: func(p *Parser) {
				p.Scan()
				p.Camouflage("test")
				p.Peek()
			},
			wantPanic: true,
		},
		{
			name:  "next and peek",
			input: "$#abc",
			do: func(p *Parser) {
				p.Scan()
				p.Next()
				p.Peek()
			},
			wantScan:  scanner.Ident,
			wantToken: "abc",
		},
		{
			name:  "scan token camouflage token",
			input: "abc",
			do: func(p *Parser) {
				p.Scan()
				p.TokenText()
				p.Camouflage("test")
			},
			wantScan:  scanner.Ident,
			wantToken: "abc",
		},
		{
			name:  "scan token peek camouflage token",
			input: "abc",
			do: func(p *Parser) {
				p.Scan()
				p.TokenText()
				p.Peek()
				p.Camouflage("test")
			},
			wantScan:  scanner.Ident,
			wantToken: "abc",
		},
		{
			name:  "tokenize all whitespace",
			input: "foo\tbar\nbaz",
			do: func(p *Parser) {
				p.SetWhitespace()
				p.Scan()
			},
			wantScan:  '\t',
			wantToken: "\t",
		},
		{
			name:  "custom ident",
			input: "$#foo",
			do: func(p *Parser) {
				p.SetIsIdentRuneFunc(func(ch rune, i int) bool { return unicode.IsLetter(ch) || ch == '#' })
				p.Scan()
			},
			wantScan:  scanner.Ident,
			wantToken: "#foo",
		},
		{
			name:  "do not scan idents",
			input: "abc",
			do: func(p *Parser) {
				p.SetMode(scanner.GoTokens ^ scanner.ScanIdents)
				p.Scan()
			},
			wantScan:  'b',
			wantToken: "b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				err := recover()
				if err != nil && !tt.wantPanic {
					t.Fatalf("unexpected panic: %v", err)
				}
			}()

			p := newParser(tt.input, tt.Language)
			tt.do(p)
			if tt.wantPanic {
				return
			}
			scan := p.Scan()
			token := p.TokenText()

			if scan != tt.wantScan || token != tt.wantToken {
				t.Errorf("Parser.Scan() = %v (%v), want %v (%v)", scan, token, tt.wantScan, tt.wantToken)
			}
		})
	}
}
