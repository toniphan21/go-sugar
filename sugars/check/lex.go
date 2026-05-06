package check

import (
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/ebnf"
)

// ---
// StatementBoundary = ";" | "{" .
// Check = StatementBoundary CheckKeyword OperandName .
// CheckKeyword = "check" .
// ---

type state int

var states = struct {
	Start          state
	Running        state
	UseKeyword     state
	UseOperandName state
	End            state
}{
	Start:      state(0),
	Running:    state(1),
	UseKeyword: state(2),
	End:        state(4),
}

func newRecognizer() sugar.LexicalParser {
	see := &sugar.LexemePredicate{}
	do := &nodeBuilder{node: &node{}}

	table := sugar.NewTransitionTable[state]()

	table.
		Add(states.Start, see.StatementBoundary, states.Running, do.begin).
		Add(states.Start, see.Any, states.Start, do.failed)

	table.
		Use(states.Running, KeywordParser(), sugar.TransitionControl[state]{
			FirstTake:     do.collectPos,
			SuccessMoveTo: states.UseKeyword,
			SuccessAction: invoke(do.collectKeyword),
			ErrorMoveTo:   states.Start,
			ErrorAction:   use(do.failed),
		})

	table.
		Use(states.UseKeyword, ebnf.OperandNameParser(), sugar.TransitionControl[state]{
			SuccessMoveTo: states.End,
			SuccessAction: invoke(do.collectOperandName),
			ErrorMoveTo:   states.Start,
			ErrorAction:   use(do.failed),
		})

	return sugar.NewLexicalParser(table, states.Start, states.End, do)
}
