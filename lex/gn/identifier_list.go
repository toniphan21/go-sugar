package gn

import (
	"nhatp.com/go/sugar"
)

// ---
// see: https://go.dev/ref/spec#IdentifierList
//
// IdentifierList = identifier { "," identifier } .
// ---

/*
Diagram(
  Start({type:'complex'}),
  Stack('ident'),
  ZeroOrMore(Sequence(',', 'ident')),
  End({type:'complex'})
)
*/

type IdentifierList struct {
	Identifiers []string
}

const IdentifierListParserID = "lex/gn.IdentifierList"

func IdentifierListParser() sugar.LexicalParser {
	const start, running, expectIdent, end = "start", "running", "expect-ident", "end"
	see := &sugar.LexemePredicate{}
	builder := sugar.NewNodeBuilder[IdentifierList]()

	doFail := builder.Fail
	doCollectIdent := builder.Collect("Ident", func(n *IdentifierList, l sugar.Lexeme) {
		n.Identifiers = append(n.Identifiers, l.Lit)
	})

	table := sugar.NewTransitionTable[string](IdentifierListParserID).
		Add(start, see.Ident, running, doCollectIdent).
		Add(start, see.IsNotIdent, end, doFail).
		Add(running, see.Comma, expectIdent).
		Peek(running, see.Any, end).
		Add(expectIdent, see.Ident, running, doCollectIdent).
		Add(expectIdent, see.IsNotIdent, end, doFail)

	return sugar.NewLexicalParser(IdentifierListParserID, table, start, end, builder)
}
