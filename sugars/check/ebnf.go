package check

import "nhatp.com/go/sugar"

type Keyword struct {
	Lex sugar.Lexeme
}

type keywordState int

var keywordStates = struct {
	Start keywordState
	End   keywordState
}{
	Start: keywordState(0),
	End:   keywordState(1),
}

func KeywordParser() sugar.LexicalParser {
	see := &sugar.LexemePredicate{}
	do := &keywordBuilder{}
	state := keywordStates

	table := sugar.NewTransitionTable[keywordState]()
	table.
		Add(state.Start, see.IdentMatch("check"), state.End, do.collect).
		Add(state.Start, see.Any, state.End, do.failed)

	return sugar.NewLexicalParser(table, state.Start, state.End, do)
}

type keywordBuilder struct {
	lex   *sugar.Lexeme
	error bool
}

func (b *keywordBuilder) Reset() {
	b.lex = nil
	b.error = false
}

func (b *keywordBuilder) Build() (Keyword, bool) {
	if b.error || b.lex == nil {
		return Keyword{}, false
	}
	return Keyword{Lex: *b.lex}, true
}

func (b *keywordBuilder) collect(lex sugar.Lexeme) {
	b.lex = &lex
}

func (b *keywordBuilder) failed(lex sugar.Lexeme) {
	b.error = true
}

var _ sugar.LexicalNodeBuilder[Keyword] = (*keywordBuilder)(nil)
