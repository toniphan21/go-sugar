package lex

import (
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/lex/gn"
)

// ---
// IdentifierLHS = IdentifierList ":=" .
// ---

type IdentifierLHS struct {
	Identifiers []string
}

const IdentifierLHSParserID = "lex.IdentifierLHS"

func IdentifierLHSParser() sugar.LexicalParser {
	const start, end = "start", "end"
	see := &sugar.LexemePredicate{}
	builder := sugar.NewNodeBuilder[IdentifierLHS]()

	doFail := builder.Fail
	doCollect := builder.CollectInner("IdentifierList", func(n *IdentifierLHS, d any, l sugar.Lexeme) {
		il, ok := d.(gn.IdentifierList)
		if !ok {
			builder.Fail(l)
			return
		}
		n.Identifiers = il.Identifiers
	})

	table := sugar.NewTransitionTable[string]()

	table.Use(start, gn.IdentifierListParser(), sugar.TransitionControl[string]{
		WhenSuccess: func(_ sugar.LexicalParser, data any, lex sugar.Lexeme) string {
			if see.Define(lex) {
				doCollect(data, lex)
			} else {
				doFail(lex)
			}
			return end
		},
	})

	return sugar.NewLexicalParser(IdentifierLHSParserID, table, start, end, builder)
}
