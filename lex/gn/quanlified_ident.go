package gn

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

/*
Diagram(
  Start({type:'complex'}),
  Stack('ident'),
  Optional(Sequence('.', 'ident')),
  End({type:'complex'})
)
*/

type TypeName = qualifiedIdent

const TypeNameParserID = "lex/gn.TypeName"

func TypeNameParser() sugar.LexicalParser {
	return newQualifiedIdentParser(TypeNameParserID)
}

type OperandName = qualifiedIdent

const OperandNameParserID = "lex/gn.OperandName"

func OperandNameParser() sugar.LexicalParser {
	return newQualifiedIdentParser(OperandNameParserID)
}

type qualifiedIdent struct {
	PackageName *string
	Identifier  string
	current     *string
}

func newQualifiedIdentParser(name string) sugar.LexicalParser {
	const start, running, expectIdent, end = "start", "running", "expect-ident", "end"
	see := &sugar.LexemePredicate{}
	builder := sugar.NewNodeBuilder[qualifiedIdent]()

	doFail := builder.Fail
	doCollect := builder.Collect("-", func(n *qualifiedIdent, l sugar.Lexeme) {
		if n.current == nil {
			n.Identifier = l.Lit
			n.current = new(l.Lit)
		} else {
			tmp := n.current
			n.PackageName = tmp
			n.Identifier = l.Lit
		}
	})

	table := sugar.NewTransitionTable[string](name)

	table.
		Add(start, see.Ident, running, doCollect).
		Add(start, see.IsNotIdent, end, doFail)

	table.
		Add(running, see.Period, expectIdent).
		Add(running, see.Any, end)

	table.
		Add(expectIdent, see.Ident, running, doCollect).
		Add(expectIdent, see.IsNotIdent, end, doFail)

	return sugar.NewLexicalParser(name, table, start, end, builder)
}
