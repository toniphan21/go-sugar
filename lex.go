package sugar

import (
	"go/scanner"
	"go/token"
)

type Lexeme struct {
	Tok    token.Token
	Lit    string
	Line   int
	Column int
	Offset int
}

func Lex(content []byte) []Lexeme {
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(content))

	var s scanner.Scanner
	s.Init(file, content, nil, 0)

	var result []Lexeme
	for {
		p, t, l := s.Scan()
		if t == token.EOF {
			break
		}

		pos := fset.Position(p)
		lex := Lexeme{
			Tok:    t,
			Lit:    l,
			Line:   pos.Line,
			Column: pos.Column,
			Offset: pos.Offset,
		}
		if lex.Lit == "" {
			lex.Lit = t.String()
		}
		result = append(result, lex)
	}
	return result
}
