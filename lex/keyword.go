package lex

import "nhatp.com/go/sugar"

type Keyword struct {
	Pos sugar.Lexeme
}

const KeywordParserID = "lex.Keyword"

func KeywordParser(keyword string) sugar.LexicalParser {
	const start, end = "start", "end"
	see := &sugar.LexemePredicate{}
	builder := sugar.NewNodeBuilder[Keyword]()

	doFail := builder.Fail
	doCollectPos := builder.Collect("Pos", func(n *Keyword, lex sugar.Lexeme) {
		n.Pos = lex
	})

	table := sugar.NewTransitionTable[string]().
		Add(start, see.IdentMatch(keyword), end, doCollectPos).
		Add(start, see.Any, end, doFail)

	return sugar.NewLexicalParser(KeywordParserID, table, start, end, builder)
}
