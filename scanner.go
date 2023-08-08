package gval

import (
	"io"
	"strconv"
	"text/scanner"
)

type defaultScanner struct {
	scanner.Scanner
}

type Scanner interface {
	Init(reader io.Reader)
	SetError(func(Scanner, string))
	SetFilename(string)
	SetWhitespace(uint64)
	GetWhitespace() uint64
	SetMode(uint)
	GetMode() uint
	SetIsIdentRune(func(rune, int) bool)
	GetIsIdentRune() func(rune, int) bool
	Scan() rune
	Peek() rune
	Next() rune
	TokenText() string
	Pos() scanner.Position
	GetPosition() scanner.Position
	Unquote(string) (string, error)
}

func (s *defaultScanner) Init(reader io.Reader) {
	s.Scanner.Init(reader)
}

func (s *defaultScanner) SetError(fn func(s Scanner, msg string)) {
	s.Scanner.Error = func(_ *scanner.Scanner, msg string) {
		fn(s, msg)
	}
}

func (s *defaultScanner) SetFilename(filename string) {
	s.Scanner.Filename = filename
}

func (s *defaultScanner) SetWhitespace(ws uint64) {
	s.Scanner.Whitespace = ws
}

func (s *defaultScanner) GetWhitespace() uint64 {
	return s.Scanner.Whitespace
}

func (s *defaultScanner) SetMode(m uint) {
	s.Scanner.Mode = m
}

func (s *defaultScanner) GetMode() uint {
	return s.Scanner.Mode
}

func (s *defaultScanner) SetIsIdentRune(fn func(ch rune, i int) bool) {
	s.Scanner.IsIdentRune = fn
}

func (s *defaultScanner) GetIsIdentRune() func(ch rune, i int) bool {
	return s.Scanner.IsIdentRune
}

func (s *defaultScanner) GetPosition() scanner.Position {
	return s.Scanner.Position
}

func (*defaultScanner) Unquote(s string) (string, error) {
	return strconv.Unquote(s)
}
