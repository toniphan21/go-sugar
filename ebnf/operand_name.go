package ebnf

import "nhatp.com/go/sugar"

type OperandName = qualifiedIdent

func OperandNameParser() sugar.LexicalParser {
	state := qualifiedIdentStates
	builder := &qualifiedIdentBuilder{}

	return sugar.NewLexicalParser[qualifiedIdentState, *qualifiedIdentBuilder, OperandName](
		qualifiedIdentTransitionTable(builder),
		state.Start,
		state.End,
		builder,
	)
}
