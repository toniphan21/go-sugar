package check

import (
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/lex"
)

// ---
// StatementBoundary = ";" | "{" .
// IdentifierLHS = IdentifierList ":=" .
// CallExpr = SelectorPath CallSuffix .
//
// Check = StatementBoundary [ IdentifierLHS ] "check" CallExpr .
// ---

type Statement struct {
	isCompleted bool
	pos         *sugar.Lexeme
	end         *sugar.Lexeme
	identifiers []string
}

const LexicalParserID = "sugars/check.LexicalParser"

func LexicalParser() sugar.LexicalParser {
	const start, running, expectCheck, afterCheck, end = "start", "running", "expect-check", "after-check", "end"
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
	doCollectAfterCheck := builder.CollectInner("AfterCheck", func(n *Statement, d any, l sugar.Lexeme) {
		sugar.CollectBuilderDataOrFail(builder, d, l, func(v lex.CallExpr) {
			n.identifiers = v.Identifiers
			n.end = &v.CallEnd
		})
	})

	keywordParser := lex.KeywordParser("check")
	identifierLHSParser := lex.IdentifierLHSParser()

	table := sugar.NewTransitionTable[string]()

	table.
		Add(start, see.StatementBoundary, running, doBegin).
		Add(start, see.Any, start, doFail)

	table.
		Longest(running, sugar.TransitionControl[string]{
			FirstTake:   doCollectPos,
			ErrorMoveTo: start,
			ErrorAction: doPropagateFail,
			WhenSuccess: func(p sugar.LexicalParser, d any, l sugar.Lexeme) (string, int) {
				switch {
				case p.Is(keywordParser):
					return afterCheck, 0

				case p.Is(identifierLHSParser):
					return expectCheck, 0

				default:
					doFail(l)
					return start, 0
				}
			},
		}, keywordParser, identifierLHSParser)

	table.
		Use(expectCheck, keywordParser, sugar.TransitionControl[string]{
			SuccessMoveTo: afterCheck,
			ErrorMoveTo:   start,
			ErrorAction:   doPropagateFail,
		})

	table.
		Use(afterCheck, lex.CallExprParser(), sugar.TransitionControl[string]{
			SuccessMoveTo:  end,
			SuccessAction:  doCollectAfterCheck,
			SuccessPutBack: 1, // CallExpr ends With StatementBoundary which is a start point so we need to putback
			ErrorMoveTo:    start,
			ErrorAction:    doPropagateFail,
		})

	return sugar.NewLexicalParser(LexicalParserID, table, start, end, builder)
}
