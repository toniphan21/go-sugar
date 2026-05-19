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

/*
Diagram(
  Start({type:'complex'}),
  Stack('identifier'),
  Stack(":="),
  End({type:'complex'})
)
*/

const IdentifierLHSParserID = "lex.IdentifierLHS"

func IdentifierLHSParser() sugar.LexicalParser {
	const start, expectDefine, end = "start", "expect-define", "end"
	see := &sugar.LexemePredicate{}
	builder := sugar.NewNodeBuilder[IdentifierLHS]()

	doFail := builder.Fail
	doPropagateFail := builder.FailInner
	doCollect := builder.CollectInner("IdentifierList", func(n *IdentifierLHS, d any, l sugar.Lexeme) {
		sugar.CollectBuilderDataOrFail(builder, d, l, func(v gn.IdentifierList) {
			n.Identifiers = v.Identifiers
		})
	})

	table := sugar.NewTransitionTable[string](IdentifierLHSParserID).
		Use(start, gn.IdentifierListParser(), sugar.TransitionControl[string]{
			ErrorMoveTo:   end,
			ErrorAction:   doPropagateFail,
			SuccessMoveTo: expectDefine,
			SuccessAction: doCollect,
		}).
		Add(expectDefine, see.Define, end).
		Add(expectDefine, see.Any, end, doFail)

	return sugar.NewLexicalParser(IdentifierLHSParserID, table, start, end, builder)
}
