package ebnf

import "nhatp.com/go/sugar"

type TypeName = qualifiedIdent

func TypeNameParser() sugar.LexicalParser {
	state := qualifiedIdentStates
	builder := &qualifiedIdentBuilder{}

	return sugar.NewLexicalParser[qualifiedIdentState, *qualifiedIdentBuilder, TypeName](
		qualifiedIdentTransitionTable(builder),
		state.Start,
		state.End,
		builder,
	)
}
