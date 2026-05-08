package check

import (
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/lex"
	"nhatp.com/go/sugar/lex/gn"
)

// ---
// StatementBoundary = ";" | "{" .
// IdentifierLHS = IdentifierList ":=" .
//
// Check = StatementBoundary [ IdentifierLHS ] "check" OperandName .
// ---

type Statement struct {
	isCompleted bool
	pos         *sugar.Lexeme
	end         *sugar.Lexeme
	operandPkg  *string
	operandName string
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
	doCollectOperandName := builder.CollectInner("OperandName", func(n *Statement, d any, l sugar.Lexeme) {
		if data, ok := d.(gn.OperandName); ok {
			n.operandPkg = data.PackageName
			n.operandName = data.Identifier
			n.end = &l
		}
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
			WhenSuccess: func(p sugar.LexicalParser, d any, l sugar.Lexeme) string {
				switch {
				case p.Is(keywordParser):
					return afterCheck

				case p.Is(identifierLHSParser):
					return expectCheck

				default:
					doFail(l)
					return start
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
		Use(afterCheck, gn.OperandNameParser(), sugar.TransitionControl[string]{
			SuccessMoveTo: end,
			SuccessAction: doCollectOperandName,
			ErrorMoveTo:   start,
			ErrorAction:   doPropagateFail,
		})

	return sugar.NewLexicalParser(LexicalParserID, table, start, end, builder).Debug()
}
