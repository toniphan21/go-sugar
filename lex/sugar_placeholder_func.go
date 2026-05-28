package lex

import (
	"fmt"

	"nhatp.com/go/sugar"
)

/*
Diagram(
  Start({type:'complex'}),
  Stack('__sugar_[keyword]__'),
  Stack('('),
  Stack('*'),
  ZeroOrMore('*'),
  Stack(')'),
  Stack(';'),
  End({type:'complex'})
)
*/

/*
Diagram(
  Start({type:'complex'}),
  Stack(NonTerminal('doCollectPos'), '__sugar_[keyword]__'),
  Stack('('),
  Stack(NonTerminal('doCollectInnerPos'), '*'),
  ZeroOrMore('*'),
  Stack(NonTerminal('doCollectInnerEnd'), ')'),
  Stack(NonTerminal('doCollectEnd'), ';'),
  End({type:'complex'})
)
*/

type SugarPlaceholderFunc struct {
	pos      sugar.Lexeme
	end      sugar.Lexeme
	innerPos sugar.Lexeme
	innerEnd sugar.Lexeme
	keyword  string
}

func (n SugarPlaceholderFunc) Pos() sugar.Lexeme {
	return n.pos
}

func (n SugarPlaceholderFunc) End() sugar.Lexeme {
	return n.end
}

func (n SugarPlaceholderFunc) InnerPos() sugar.Lexeme {
	return n.innerPos
}

func (n SugarPlaceholderFunc) InnerEnd() sugar.Lexeme {
	return n.innerEnd
}

func (n SugarPlaceholderFunc) Keyword() string {
	return n.keyword
}

const SugarPlaceholderFuncID = "lex.SugarPlaceholderFunc"

func SugarPlaceholderFuncName(keyword string) string {
	return fmt.Sprintf("__sugar_%s__", keyword)
}

func SugarPlaceholderFuncParser(keyword string) sugar.LexicalParser {
	const start, expectLParen, expectAny, running, end = "start", "expect-lparen", "expect-any", "running", "end"
	see := &sugar.LexemePredicate{}
	builder := sugar.NewNodeBuilder[SugarPlaceholderFunc]()

	const deep = "deep"
	doFail := builder.Fail
	doBegin := builder.Collect("begin", func(n *SugarPlaceholderFunc, l sugar.Lexeme) {
		builder.Error = false
		n.keyword = keyword
		n.pos = l
	})
	doCollectInnerPos := builder.Collect("inner-pos", func(n *SugarPlaceholderFunc, l sugar.Lexeme) {
		n.innerPos = l
	})
	doIncDeep := builder.Collect("inc", func(n *SugarPlaceholderFunc, l sugar.Lexeme) {
		builder.CounterInc(deep)
	})
	doDecDeep := builder.Collect("inc", func(n *SugarPlaceholderFunc, l sugar.Lexeme) {
		builder.CounterDec(deep)
		n.innerEnd = l
	})
	doCollect := builder.Collect("end", func(n *SugarPlaceholderFunc, l sugar.Lexeme) {
		n.end = l
	})
	doAtStatementBoundary := func(lex sugar.Lexeme) {
		if builder.Counter(deep) == 0 {
			doCollect(lex)
		} else {
			doFail(lex)
		}
	}

	table := sugar.NewTransitionTable[string](CallExprID)
	table.
		Add(start, see.IdentMatch(SugarPlaceholderFuncName(keyword)), expectLParen, doBegin).
		Add(start, see.Any, end, doFail).
		Add(expectLParen, see.LeftParen, expectAny, doIncDeep).
		Add(expectLParen, see.Any, end, doFail).
		Add(expectAny, see.LeftParen, running, doIncDeep, doCollectInnerPos).
		Add(expectAny, see.RightParen, end, doFail).
		Add(expectAny, see.StatementBoundary, end, doFail).
		Add(expectAny, see.Any, running, doCollectInnerPos).
		Add(running, see.LeftParen, running, doIncDeep).
		Add(running, see.RightParen, running, doDecDeep).
		Add(running, see.StatementBoundary, end, doAtStatementBoundary).
		Add(running, see.Any, running)

	return sugar.NewLexicalParser(SugarPlaceholderFuncID, table, start, end, builder)
}
