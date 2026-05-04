package parse

import (
	"nhatp.com/go/sugar"
)

/* see: https://go.dev/ref/spec#OperandName

EBNF:
	OperandName = identifier | QualifiedIdent .
	QualifiedIdent = PackageName "." identifier .
	PackageName   = identifier .
*/

type OperandName struct {
	PackageName *string
	Identifier  string
}

type operandNameState int

const (
	operandNameStateStart operandNameState = iota
	operandNameStateRunning
	operandNameStateExpectIdent
	operandNameStateEnd
)

func OperandNameLexicalParser() LexicalParser[*OperandName] {
	do := &operandNameBuilder{}
	on := &sugar.LexemePredicate{}

	transitions := []sugar.Transition[operandNameState]{
		{From: operandNameStateStart, Event: on.Ident, To: operandNameStateRunning, Action: do.collect},
		{From: operandNameStateStart, Event: on.IsNotIdent, To: operandNameStateEnd, Action: do.failed},
		{From: operandNameStateRunning, Event: on.Period, To: operandNameStateExpectIdent},
		{From: operandNameStateRunning, Event: on.Any, To: operandNameStateEnd},
		{From: operandNameStateExpectIdent, Event: on.Ident, To: operandNameStateRunning, Action: do.collect},
		{From: operandNameStateExpectIdent, Event: on.IsNotIdent, To: operandNameStateEnd, Action: do.failed},
	}

	return &lexicalParserImpl[operandNameState, *operandNameBuilder, *OperandName]{
		transitions:  transitions,
		builder:      do,
		current:      operandNameStateStart,
		initialState: operandNameStateStart,
		statusMapper: func(state operandNameState, builder *operandNameBuilder) sugar.Status {
			switch state {
			case operandNameStateEnd:
				if builder.error {
					return sugar.StatusFailed
				}
				return sugar.StatusCompleted

			default:
				return sugar.StatusRunning
			}
		},
	}
}

type operandNameBuilder struct {
	pkgName *string
	ident   *string
	error   bool
}

func (b *operandNameBuilder) reset() {
	b.pkgName = nil
	b.ident = nil
	b.error = false
}

func (b *operandNameBuilder) build() *OperandName {
	if b.error {
		return nil
	}
	return &OperandName{PackageName: b.pkgName, Identifier: *b.ident}
}

func (b *operandNameBuilder) collect(lex sugar.Lexeme) {
	if b.ident != nil {
		b.pkgName = b.ident
		b.ident = new(lex.Lit)
	} else {
		b.ident = new(lex.Lit)
	}
}

func (b *operandNameBuilder) failed(lex sugar.Lexeme) {
	b.error = true
}
