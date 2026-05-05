package ebnf

import (
	"nhatp.com/go/sugar"
)

// ---
// based type of OperandName and TypeName. See:
//  - https://go.dev/ref/spec#OperandName
//  - https://go.dev/ref/spec#Type
//
// OperandName = identifier | QualifiedIdent .
// TypeName = identifier | QualifiedIdent .
//
// QualifiedIdent = PackageName "." identifier .
// PackageName   = identifier .
// ---

type qualifiedIdent struct {
	PackageName *string
	Identifier  string
}

type qualifiedIdentState int

var qualifiedIdentStates = struct {
	Start       qualifiedIdentState
	Running     qualifiedIdentState
	ExpectIdent qualifiedIdentState
	End         qualifiedIdentState
}{
	Start:       qualifiedIdentState(0),
	Running:     qualifiedIdentState(1),
	ExpectIdent: qualifiedIdentState(2),
	End:         qualifiedIdentState(3),
}

func qualifiedIdentTransitionTable(do *qualifiedIdentBuilder) sugar.TransitionTable[qualifiedIdentState] {
	see := &sugar.LexemePredicate{}
	state := qualifiedIdentStates

	table := sugar.NewTransitionTable[qualifiedIdentState]()

	table.
		Add(state.Start, see.Ident, state.Running, do.collect).
		Add(state.Start, see.IsNotIdent, state.End, do.failed)

	table.
		Add(state.Running, see.Period, state.ExpectIdent).
		Add(state.Running, see.Any, state.End)

	table.
		Add(state.ExpectIdent, see.Ident, state.Running, do.collect).
		Add(state.ExpectIdent, see.IsNotIdent, state.End, do.failed)

	return table
}

type qualifiedIdentBuilder struct {
	pkgName *string
	ident   *string
	error   bool
}

func (b *qualifiedIdentBuilder) Reset() {
	b.pkgName = nil
	b.ident = nil
	b.error = false
}

func (b *qualifiedIdentBuilder) Build() (qualifiedIdent, bool) {
	if b.error {
		return qualifiedIdent{}, false
	}
	return qualifiedIdent{PackageName: b.pkgName, Identifier: *b.ident}, true
}

func (b *qualifiedIdentBuilder) failed(lex sugar.Lexeme) {
	b.error = true
}

func (b *qualifiedIdentBuilder) collect(lex sugar.Lexeme) {
	if b.ident != nil {
		b.pkgName = b.ident
		b.ident = new(lex.Lit)
	} else {
		b.ident = new(lex.Lit)
	}
}
