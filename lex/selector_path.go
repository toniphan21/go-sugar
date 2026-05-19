package lex

import "nhatp.com/go/sugar"

// ---
// SelectorPath = identifier { "." identifier } .
// ---

/*
Diagram(
  Start({type:'complex'}),
  Stack('ident'),
  ZeroOrMore(Sequence('.', 'ident')),
  End({type:'complex'})
)
*/

type SelectorPath struct {
	Identifiers []string
}

const SelectorPathID = "lex.SelectorPath"

func SelectorPathParser() sugar.LexicalParser {
	const start, expectDot, expectIdent, end = "start", "expect-dot", "expect-ident", "end"
	see := &sugar.LexemePredicate{}
	builder := sugar.NewNodeBuilder[SelectorPath]()

	doFail := builder.Fail
	doCollect := builder.Collect("-", func(n *SelectorPath, l sugar.Lexeme) {
		n.Identifiers = append(n.Identifiers, l.Lit)
	})

	table := sugar.NewTransitionTable[string](SelectorPathID)

	table.
		Add(start, see.Ident, expectDot, doCollect).
		Add(start, see.Any, end, doFail).
		Add(expectDot, see.Period, expectIdent).
		Peek(expectDot, see.Any, end).
		Add(expectIdent, see.Ident, expectDot, doCollect).
		Add(expectIdent, see.Any, end, doFail)

	return sugar.NewLexicalParser(SelectorPathID, table, start, end, builder)
}
