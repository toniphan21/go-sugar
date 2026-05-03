package sugar

import (
	"go/scanner"
	"go/token"
)

type Lexeme interface {
	Tok() token.Token
	Lit() string
	Raw() string
	Offset() int
	Line() int
	Column() int
}

type Token struct {
	tok    token.Token
	lit    string
	raw    string
	line   int
	column int
	offset int
}

func (t *Token) Tok() token.Token {
	return t.tok
}

func (t *Token) Lit() string {
	return t.lit
}

func (t *Token) Raw() string {
	return t.raw
}

func (t *Token) SetRaw(raw string) {
	t.raw = raw
}

func (t *Token) Line() int {
	return t.line
}

func (t *Token) Column() int {
	return t.column
}

func (t *Token) Offset() int {
	return t.offset
}

var _ Lexeme = (*Token)(nil)

func NewToken(tok token.Token, lit string, line, column, offset int) *Token {
	return &Token{
		tok:    tok,
		lit:    lit,
		line:   line,
		column: column,
		offset: offset,
	}
}

type ScanFunc func(prev Lexeme, current *Token, next *Token, source []byte) Lexeme

func Scan(content []byte, scan ScanFunc) []Lexeme {
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(content))

	var s scanner.Scanner
	s.Init(file, content, nil, 0)

	var tokens []Lexeme
	var current *Token

	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}

		position := fset.Position(pos)
		next := &Token{tok: tok, lit: lit, line: position.Line, column: position.Column, offset: position.Offset}
		if next.lit == "" {
			next.lit = tok.String()
		}

		if current != nil {
			current.raw = string(content[current.offset:position.Offset])
			var lex Lexeme = current
			if scan != nil {
				var prev Lexeme
				if len(tokens) > 0 {
					prev = tokens[len(tokens)-1]
				}
				lex = scan(prev, current, next, content)
			}
			tokens = append(tokens, lex)
		}

		current = next
	}

	if current != nil {
		current.raw = string(content[current.offset:])
		var lex Lexeme = current
		if scan != nil {
			var prev Lexeme
			if len(tokens) > 0 {
				prev = tokens[len(tokens)-1]
			}
			lex = scan(prev, current, nil, content)
		}
		tokens = append(tokens, lex)
	}
	return tokens
}
