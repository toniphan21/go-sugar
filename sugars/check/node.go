package check

import "nhatp.com/go/sugar"

type node struct {
}

type nodeBuilder struct {
}

func (b *nodeBuilder) reset(lex sugar.Lexeme) {
}

func (b *nodeBuilder) appendVariable(lex sugar.Lexeme) {
}

func (b *nodeBuilder) setOpAssign(lex sugar.Lexeme) {
}

func (b *nodeBuilder) setOpDefine(lex sugar.Lexeme) {
}

func (b *nodeBuilder) appendOperand(lex sugar.Lexeme) {
}

func (b *nodeBuilder) clearLHS(lex sugar.Lexeme) {
}

func (b *nodeBuilder) build() node {
	return node{}
}
