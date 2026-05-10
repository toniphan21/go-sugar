package lex

import (
	"nhatp.com/go/sugar"
)

// ---
// CallExpr = SelectorPath CallSuffix .
// ---

type CallExpr struct {
	Identifiers []string
	CallPos     sugar.Lexeme
	CallEnd     sugar.Lexeme
}

const CallExprID = "lex.CallExpr"

func CallExprParser() sugar.LexicalParser {
	const start, running, end = "start", "running", "end"
	builder := sugar.NewNodeBuilder[CallExpr]()

	doPropagateFail := builder.FailInner
	doCollectSelectorPath := builder.CollectInner("SelectorPath", func(n *CallExpr, d any, l sugar.Lexeme) {
		builder.Error = false
		sugar.CollectBuilderDataOrFail(builder, d, l, func(v SelectorPath) {
			n.Identifiers = make([]string, len(v.Identifiers))
			copy(n.Identifiers, v.Identifiers)
		})
	})
	doCollectCallSuffix := builder.CollectInner("CallSuffix", func(n *CallExpr, d any, l sugar.Lexeme) {
		sugar.CollectBuilderDataOrFail(builder, d, l, func(v CallSuffix) {
			n.CallPos = v.Pos
			n.CallEnd = v.End
		})
	})

	table := sugar.NewTransitionTable[string]()

	table.
		Use(start, SelectorPathParser(), sugar.TransitionControl[string]{
			SuccessMoveTo:  running,
			SuccessAction:  doCollectSelectorPath,
			SuccessPutBack: 1, // SelectorPath ends with !Ident, so we putback 1 for CallSuffix to start correctly
			ErrorMoveTo:    end,
			ErrorAction:    doPropagateFail,
		}).
		Use(running, CallSuffixParser(), sugar.TransitionControl[string]{
			SuccessMoveTo: end,
			SuccessAction: doCollectCallSuffix,
			ErrorMoveTo:   end,
			ErrorAction:   doPropagateFail,
		})

	return sugar.NewLexicalParser(CallExprID, table, start, end, builder)
}
