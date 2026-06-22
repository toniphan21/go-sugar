package lex

import (
	"nhatp.com/go/sugar"
)

type LValue struct {
	Pos sugar.Lexeme
	End sugar.Lexeme
}

const LValueParserID = "lex.LValue"

func LValueParser() sugar.LexicalParser {

	return nil
}
