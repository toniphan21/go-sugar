package check

import (
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/ebnf"
)

type node struct {
	isCompleted bool
	pos         sugar.Lexeme
	end         sugar.Lexeme
	operandPkg  *string
	operandName string
}

func invoke[T any](fn func(data T, lex sugar.Lexeme)) func(data any, lex sugar.Lexeme) {
	return func(data any, lex sugar.Lexeme) {
		v, ok := data.(T)
		if ok {
			fn(v, lex)
		}
	}
}

func use(fn func(lex sugar.Lexeme)) func(data any, lex sugar.Lexeme) {
	return func(data any, lex sugar.Lexeme) {
		fn(lex)
	}
}

type nodeBuilder struct {
	error bool
	node  *node
}

func (b *nodeBuilder) Reset() {
	b.error = false
	b.node = new(node)
}

func (b *nodeBuilder) Build() (node, bool) {
	b.node.isCompleted = !b.error
	n := *b.node

	return n, !b.error
}

func (b *nodeBuilder) failed(lex sugar.Lexeme) {
	b.error = true
}

func (b *nodeBuilder) begin(lex sugar.Lexeme) {
	b.error = false
}

func (b *nodeBuilder) collectPos(lex sugar.Lexeme) {
	b.node.pos = lex
}

func (b *nodeBuilder) collectKeyword(data Keyword, lex sugar.Lexeme) {
}

func (b *nodeBuilder) collectOperandName(data ebnf.OperandName, lex sugar.Lexeme) {
	b.node.operandPkg = data.PackageName
	b.node.operandName = data.Identifier
	b.node.end = lex
}
