package check

import (
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/lex"
	"nhatp.com/go/sugar/lex/gn"
)

// ---
// StatementBoundary = ";" | "{" .
// Check = StatementBoundary CheckKeyword OperandName .
// CheckKeyword = "check" .
// ---

type Statement struct {
	isCompleted bool
	pos         *sugar.Lexeme
	end         *sugar.Lexeme
	operandPkg  *string
	operandName string
}

func LexicalParser() sugar.LexicalParser {
	const start, running, useKeyword, end = "start", "running", "use-keyword", "end"
	see := &sugar.LexemePredicate{}

	builder := sugar.NewNodeBuilder[Statement]().OnBuild(func(n *Statement, ok bool) {
		n.isCompleted = ok
	})

	doFail, doPropagateFail := builder.Fail, builder.FailInner
	doBegin := builder.Collect("begin", func(n *Statement, l sugar.Lexeme) {
		builder.Error = false
	})
	doCollectPos := builder.Collect("Pos", func(n *Statement, l sugar.Lexeme) {
		n.pos = &l
	})
	doCollectOperandName := builder.CollectInner("OperandName", func(n *Statement, d any, l sugar.Lexeme) {
		if data, ok := d.(gn.OperandName); ok {
			n.operandPkg = data.PackageName
			n.operandName = data.Identifier
			n.end = &l
		}
	})

	table := sugar.NewTransitionTable[string]()

	table.
		Add(start, see.StatementBoundary, running, doBegin).
		Add(start, see.Any, start, doFail)

	table.
		Use(running, lex.KeywordParser("check"), sugar.TransitionControl[string]{
			FirstTake:     doCollectPos,
			SuccessMoveTo: useKeyword,
			ErrorMoveTo:   start,
			ErrorAction:   doPropagateFail,
		})

	table.
		Use(useKeyword, gn.OperandNameParser(), sugar.TransitionControl[string]{
			SuccessMoveTo: end,
			SuccessAction: doCollectOperandName,
			ErrorMoveTo:   start,
			ErrorAction:   doPropagateFail,
		})

	return sugar.NewLexicalParser("sugars/check.LexicalParser", table, start, end, builder)
}
