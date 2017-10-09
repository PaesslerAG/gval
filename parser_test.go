package gval

import (
	"testing"
	"text/scanner"
)

func TestParser_Scan(t *testing.T) {
	tests := []struct {
		name  string
		input string
		Language
		do        func(p *Parser)
		wantScan  rune
		wantToken string
		wanPanic  bool
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
			wanPanic: true,
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
			wanPanic: true,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				err := recover()
				if err != nil && !tt.wanPanic {
					t.Fatalf("unexpected panic: %v", err)
				}
			}()

			p := newParser(tt.input, tt.Language)
			tt.do(p)
			if tt.wanPanic {
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
