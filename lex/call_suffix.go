package lex

import (
	"nhatp.com/go/sugar"
)

// ---
// StatementBoundary = ";" | "{" .
// CallSuffix = "(" TokenSeq ")" [StatementBoundary] .
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
	const terminate = "terminate"

	doFail := builder.Fail
	doBegin := builder.Collect("begin", func(n *CallSuffix, l sugar.Lexeme) {
		n.Pos = l
	})
	doIncDeep := builder.Collect("inc", func(n *CallSuffix, l sugar.Lexeme) {
		builder.CounterInc(deep)
	})
	doCollect := builder.Collect("end", func(n *CallSuffix, l sugar.Lexeme) {
		n.End = l
	})

	table := sugar.NewTransitionTable[string](CallExprID)
	table.
		Add(start, see.LeftParen, running, doBegin, doIncDeep).
		Add(start, see.Any, end, doFail).
		Route(running, func(lex sugar.Lexeme) (string, bool) {
			switch {
			case see.LeftParen(lex):
				builder.CounterInc(deep)
				return running, true

			case see.RightParen(lex):
				builder.CounterDec(deep)
				if builder.Counter(deep) == 0 {
					builder.SetFlag(terminate)
				}
				return running, true

			case see.StatementBoundary(lex):
				if builder.Counter(deep) == 0 {
					doCollect(lex)
				} else {
					doFail(lex)
				}
				return end, true

			default:
				if builder.Flag(terminate) {
					doCollect(lex)
					return end, true
				}
				return running, false
			}
		})

	return sugar.NewLexicalParser(CallSuffixID, table, start, end, builder)
}
