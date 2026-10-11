package gval

import (
	"strconv"
	"strings"
	"text/scanner"
	"unicode/utf8"
)

// scanStringLiteral scans a double-quoted string literal,
// permitting unescaped line feeds while preserving string contents and column/line tracking.
func (p *Parser) scanStringLiteral(quote rune) string {
	var b strings.Builder
	b.WriteRune(quote)
	for {
		ch := p.scanner.Next()
		if ch == scanner.EOF {
			break
		}
		b.WriteRune(ch)
		if ch == '\\' {
			esc := p.scanner.Next()
			if esc == scanner.EOF {
				break
			}
			b.WriteRune(esc)
		} else if ch == quote {
			break
		}
	}
	return b.String()
}

// unquote unquotes single-quoted, double-quoted, or raw backquoted string literals,
// permitting unescaped line feeds in double-quoted strings.
func unquote(s string) (string, error) {
	n := len(s)
	if n < 2 {
		return "", strconv.ErrSyntax
	}
	quote := s[0]
	if quote != s[n-1] {
		return "", strconv.ErrSyntax
	}
	s = s[1 : n-1]

	if quote == '`' {
		if strings.Contains(s, "`") {
			return "", strconv.ErrSyntax
		}
		if strings.Contains(s, "\r") {
			var buf []byte
			for i := 0; i < len(s); i++ {
				if s[i] != '\r' {
					buf = append(buf, s[i])
				}
			}
			return string(buf), nil
		}
		return s, nil
	}

	if quote != '"' && quote != '\'' {
		return "", strconv.ErrSyntax
	}

	if quote == '\'' {
		r, multibyte, rem, err := strconv.UnquoteChar(s, '\'')
		if err != nil || len(rem) > 0 {
			return "", strconv.ErrSyntax
		}
		if r < utf8.RuneSelf || !multibyte {
			return string([]byte{byte(r)}), nil
		}
		return string(r), nil
	}

	// Double-quoted string literal
	if !strings.Contains(s, "\\") {
		if strings.Contains(s, "\"") || !utf8.ValidString(s) {
			return "", strconv.ErrSyntax
		}
		return s, nil
	}

	var buf []byte
	for len(s) > 0 {
		r, multibyte, rem, err := strconv.UnquoteChar(s, '"')
		if err != nil {
			return "", strconv.ErrSyntax
		}
		s = rem
		if r < utf8.RuneSelf || !multibyte {
			buf = append(buf, byte(r))
		} else {
			buf = utf8.AppendRune(buf, r)
		}
	}
	return string(buf), nil
}
