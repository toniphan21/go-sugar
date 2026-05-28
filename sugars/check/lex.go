package check

import (
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/lex"
)

// ---
// IdentifierLHS = IdentifierList ":=" .
// CallExpr = SelectorPath CallSuffix .
//
// Check = [ IdentifierLHS ] "check" CallExpr .
// ---

/*
Diagram(
  Start({type:'complex'}),
  NonTerminal('doBegin'),
  Optional(Stack('IdentifierLHS', NonTerminal('doCollectIdentifiers'))),
  Comment('expect-check'),
  Stack(NonTerminal('doCollectCheckPos'), "check"),
  Comment('expect-expr'),
  Stack(NonTerminal('doCollectCheckEnd'), 'CallExpr', NonTerminal('doCollectEnd')),
  End({type:'complex'})
)
*/

type Statement struct {
	isCompleted bool
	pos         *sugar.Lexeme
	end         *sugar.Lexeme
	checkPos    *sugar.Lexeme
	checkEnd    *sugar.Lexeme
	identifiers []string
}

func (n Statement) AsSugar() (sugar.Sugar, bool) {
	if !n.isCompleted {
		return nil, false
	}

	i := &sugarImpl{
		pos:      *n.pos,
		end:      *n.end,
		checkPos: *n.checkPos,
		checkEnd: *n.checkEnd,
	}

	if len(n.identifiers) > 0 {
		i.identifiers = make([]string, len(n.identifiers))
		copy(i.identifiers, n.identifiers)
	}
	return i, true
}

const LexicalParserID = "sugars/check.LexicalParser"

func LexicalParser() sugar.LexicalParser {
	const start, expectCheck, expectExpr, end = "start", "expect-check", "expect-expr", "end"

	builder := sugar.NewNodeBuilder[Statement]().OnBuild(func(n *Statement, ok bool) {
		n.isCompleted = ok
	})

	doPropagateFail := builder.FailInner
	doBegin := builder.Collect("begin", func(n *Statement, l sugar.Lexeme) {
		builder.Error = false
		n.pos = &l
	})
	doCollectCheckPos := builder.Collect("CheckPos", func(n *Statement, l sugar.Lexeme) {
		n.checkPos = &l
	})
	doCollectCheckEnd := builder.Collect("CheckEnd", func(n *Statement, l sugar.Lexeme) {
		n.checkEnd = &l
	})
	doCollectIdentifiers := builder.CollectInner("Identifiers", func(n *Statement, d any, l sugar.Lexeme) {
		sugar.CollectBuilderDataOrFail(builder, d, l, func(v lex.IdentifierLHS) {
			n.identifiers = v.Identifiers
		})
	})
	doCollectEnd := builder.CollectInner("End", func(n *Statement, d any, l sugar.Lexeme) {
		sugar.CollectBuilderDataOrFail(builder, d, l, func(v lex.CallExpr) {
			n.end = &v.CallEnd
		})
	})

	table := sugar.NewTransitionTable[string](LexicalParserID).
		Optional(start, lex.IdentifierLHSParser(), sugar.OptionalTransitionControl[string]{
			FirstTake:     doBegin,
			MoveTo:        expectCheck,
			SuccessAction: doCollectIdentifiers,
		}).
		Use(expectCheck, lex.KeywordParser(keyword), sugar.TransitionControl[string]{
			FirstTake:     doCollectCheckPos,
			ErrorMoveTo:   start,
			ErrorAction:   doPropagateFail,
			SuccessMoveTo: expectExpr,
		}).
		Use(expectExpr, lex.CallExprParser(), sugar.TransitionControl[string]{
			FirstTake:     doCollectCheckEnd,
			SuccessMoveTo: end,
			SuccessAction: doCollectEnd,
			ErrorMoveTo:   start,
			ErrorAction:   doPropagateFail,
		})

	return sugar.NewLexicalParser(LexicalParserID, table, start, end, builder)
}
