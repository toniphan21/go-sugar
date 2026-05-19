package lex

import "nhatp.com/go/sugar"

// ---
// StatementBoundary = ";" | "{" .
// CallSuffix = "(" TokenSeq ")" StatementBoundary .
//
// /* TokenSeq: any token sequence where parentheses are balanced;
//   "(" increments depth, ")" decrements depth, CallSuffix closes when depth reaches 0 */
// ---

/*
Diagram(
  Start({type:'complex'}),
  Stack('('),
  ZeroOrMore(Sequence('*')),
  Stack(')'),
  Sequence(';'),
  End({type:'complex'})
)
*/

type CallSuffix struct {
	Pos sugar.Lexeme
	End sugar.Lexeme
}

const CallSuffixID = "lex.CallSuffix"

func CallSuffixParser() sugar.LexicalParser {
	const start, running, end = "start", "running", "end"
	see := &sugar.LexemePredicate{}
	builder := sugar.NewNodeBuilder[CallSuffix]()

	const deep = "deep"
	doFail := builder.Fail
	doBegin := builder.Collect("begin", func(n *CallSuffix, l sugar.Lexeme) {
		builder.Error = false
		n.Pos = l
	})
	doIncDeep := builder.Collect("inc", func(n *CallSuffix, l sugar.Lexeme) {
		builder.CounterInc(deep)
	})
	doDecDeep := builder.Collect("inc", func(n *CallSuffix, l sugar.Lexeme) {
		builder.CounterDec(deep)
	})
	doCollect := builder.Collect("end", func(n *CallSuffix, l sugar.Lexeme) {
		n.End = l
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
		Add(start, see.LeftParen, running, doBegin, doIncDeep).
		Add(start, see.Any, end, doFail).
		Add(running, see.LeftParen, running, doIncDeep).
		Add(running, see.RightParen, running, doDecDeep).
		Add(running, see.StatementBoundary, end, doAtStatementBoundary).
		Add(running, see.Any, running)

	return sugar.NewLexicalParser(CallSuffixID, table, start, end, builder)
}
