package sugar

import (
	"go/scanner"
	"go/token"
)

type Lexeme interface {
	Tok() token.Token
	Pos() token.Pos
	Lit() string
	Offset() int
	Raw() string
}

type Token struct {
	tok    token.Token
	pos    token.Pos
	lit    string
	offset int
	raw    string
}

func (t *Token) Tok() token.Token {
	return t.tok
}

func (t *Token) Pos() token.Pos {
	return t.pos
}

func (t *Token) Lit() string {
	return t.lit
}

func (t *Token) Offset() int {
	return t.offset
}

func (t *Token) Raw() string {
	return t.raw
}

func (t *Token) SetRaw(raw string) {
	t.raw = raw
}

var _ Lexeme = (*Token)(nil)

func NewToken(tok token.Token, pos token.Pos, offset int) *Token {
	return &Token{tok: tok, pos: pos, offset: offset}
}

type ScanFunc func(prev Lexeme, current *Token, next *Token, source []byte) Lexeme

func Scan(source []byte, scan ScanFunc) []Lexeme {
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(source))

	var s scanner.Scanner
	s.Init(file, source, nil, 0)

	var tokens []Lexeme
	var current *Token

	for {
		pos, tok, lit := s.Scan()
		offset := fset.Position(pos).Offset
		if tok == token.EOF {
			break
		}

		next := &Token{tok: tok, pos: pos, lit: lit, offset: offset}
		if next.lit == "" {
			next.lit = tok.String()
		}

		if current != nil {
			current.raw = string(source[current.offset:offset])
			var lex Lexeme = current
			if scan != nil {
				var prev Lexeme
				if len(tokens) > 0 {
					prev = tokens[len(tokens)-1]
				}
				lex = scan(prev, current, next, source)
			}
			tokens = append(tokens, lex)
		}

		current = next
	}

	if current != nil {
		current.raw = string(source[current.offset:])
		var lex Lexeme = current
		if scan != nil {
			var prev Lexeme
			if len(tokens) > 0 {
				prev = tokens[len(tokens)-1]
			}
			lex = scan(prev, current, nil, source)
		}
		tokens = append(tokens, lex)
	}
	return tokens
}
