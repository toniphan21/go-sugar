package sugar

import (
	"go/scanner"
	"go/token"
)

type Lexeme struct {
	Tok    token.Token
	Pos    token.Pos
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
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}

		position := fset.Position(pos)
		lex := Lexeme{
			Tok:    tok,
			Pos:    pos,
			Lit:    lit,
			Line:   position.Line,
			Column: position.Column,
			Offset: position.Offset,
		}
		if lex.Lit == "" {
			lex.Lit = tok.String()
		}
		result = append(result, lex)
	}
	return result
}
