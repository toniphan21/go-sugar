package ebnf

import (
	"nhatp.com/go/sugar"
)

// ---
// see: https://go.dev/ref/spec#IdentifierList
//
// IdentifierList = identifier { "," identifier } .
// ---

type IdentifierList struct {
	Identifier []string
}

type identifierListState int

var identifierListStates = struct {
	Start       identifierListState
	Running     identifierListState
	ExpectIdent identifierListState
	End         identifierListState
}{
	Start:       identifierListState(0),
	Running:     identifierListState(1),
	ExpectIdent: identifierListState(2),
	End:         identifierListState(3),
}

func IdentifierListParser() sugar.LexicalParser {
	see := &sugar.LexemePredicate{}
	do := &identifierListBuilder{}
	state := identifierListStates

	table := sugar.NewTransitionTable[identifierListState]().
		Add(state.Start, see.Ident, state.Running, do.collect).
		Add(state.Start, see.IsNotIdent, state.End, do.failed).
		Add(state.Running, see.Comma, state.ExpectIdent).
		Add(state.Running, see.Any, state.End).
		Add(state.ExpectIdent, see.Ident, state.Running, do.collect).
		Add(state.ExpectIdent, see.IsNotIdent, state.End, do.failed)

	return sugar.NewLexicalParser(table, state.Start, state.End, do)
}

type identifierListBuilder struct {
	idents []string
	error  bool
}

func (b *identifierListBuilder) Reset() {
	b.idents = nil
	b.error = false
}

func (b *identifierListBuilder) Build() (IdentifierList, bool) {
	if b.error {
		return IdentifierList{}, false
	}

	var result []string
	for _, v := range b.idents {
		result = append(result, v)
	}
	return IdentifierList{Identifier: result}, true
}

func (b *identifierListBuilder) failed(lex sugar.Lexeme) {
	b.error = true
}

func (b *identifierListBuilder) collect(lex sugar.Lexeme) {
	b.idents = append(b.idents, lex.Lit)
}

var _ sugar.LexicalNodeBuilder[IdentifierList] = (*identifierListBuilder)(nil)
