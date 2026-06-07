package require

import (
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/lex"
)

// ---
// IdentifierLHS = IdentifierList ":=" .
// CallExpr = SelectorPath CallSuffix .
// Message = literal:string .
//
// Require = [ IdentifierLHS ] "require" CallExpr [Message].
// ---

type Statement struct {
	isCompleted bool
	pos         *sugar.Lexeme
	end         *sugar.Lexeme
	requirePos  *sugar.Lexeme
	requireEnd  *sugar.Lexeme
	messagePos  *sugar.Lexeme
	messageEnd  *sugar.Lexeme
	message     *string
	identifiers []string
}

func (n Statement) AsSugar() (sugar.Sugar, bool) {
	if !n.isCompleted {
		return nil, false
	}

	i := &sugarImpl{
		pos:        *n.pos,
		end:        *n.end,
		requirePos: *n.requirePos,
		requireEnd: *n.requireEnd,
		messagePos: n.messagePos,
		messageEnd: n.messageEnd,
		message:    n.message,
	}

	if len(n.identifiers) > 0 {
		i.identifiers = make([]string, len(n.identifiers))
		copy(i.identifiers, n.identifiers)
	}
	return i, true
}

const LexicalParserID = "sugars/require.LexicalParser"

func LexicalParser() sugar.LexicalParser {
	const start, expectRequire, expectExpr, expectStatementBoundary, end = "start", "expect-require", "expect-expr", "expect-stm-boundary", "end"

	see := &sugar.LexemePredicate{}
	builder := sugar.NewNodeBuilder[Statement]().OnBuild(func(n *Statement, ok bool) {
		n.isCompleted = ok
	})

	doPropagateFail := builder.FailInner
	doFail := builder.Fail
	doBegin := builder.Collect("begin", func(n *Statement, l sugar.Lexeme) {
		builder.Error = false
		n.pos = &l
	})
	doCollectCheckPos := builder.Collect("RequirePos", func(n *Statement, l sugar.Lexeme) {
		n.requirePos = &l
	})
	doCollectCheckEnd := builder.Collect("RequireEnd", func(n *Statement, l sugar.Lexeme) {
		n.requireEnd = &l
	})
	doCollectIdentifiers := builder.CollectInner("Identifiers", func(n *Statement, d any, l sugar.Lexeme) {
		sugar.CollectBuilderDataOrFail(builder, d, l, func(v lex.IdentifierLHS) {
			n.identifiers = v.Identifiers
		})
	})
	doCollectEnd := builder.Collect("End", func(n *Statement, l sugar.Lexeme) {
		n.messageEnd = &l
		n.end = &l
	})

	table := sugar.NewTransitionTable[string](LexicalParserID).
		Optional(start, lex.IdentifierLHSParser(), sugar.OptionalTransitionControl[string]{
			FirstTake:     doBegin,
			MoveTo:        expectRequire,
			SuccessAction: doCollectIdentifiers,
		}).
		Use(expectRequire, lex.KeywordParser(keyword), sugar.TransitionControl[string]{
			FirstTake:     doCollectCheckPos,
			ErrorMoveTo:   start,
			ErrorAction:   doPropagateFail,
			SuccessMoveTo: expectExpr,
		}).
		Use(expectExpr, lex.CallExprParser(), sugar.TransitionControl[string]{
			FirstTake: doCollectCheckEnd,
			WhenSuccess: func(p sugar.LexicalParser, d any, l sugar.Lexeme) (string, int) {
				if v, ok := d.(lex.CallExpr); ok {
					if see.String(l) {
						builder.Node.messagePos = &l
						builder.Node.message = new(l.Lit)

						return expectStatementBoundary, 0
					}

					builder.Node.end = &v.CallEnd
					return end, 0
				}

				builder.Fail(l)
				return end, 0
			},
			ErrorMoveTo: start,
			ErrorAction: doPropagateFail,
		}).
		Add(expectStatementBoundary, see.StatementBoundary, end, doCollectEnd).
		Add(expectStatementBoundary, see.Any, end, doFail)

	return sugar.NewLexicalParser(LexicalParserID, table, start, end, builder)
}
