package check

import "nhatp.com/go/sugar"

type state int

const (
	stateIdle state = iota
	stateStart
	stateTarget
	stateExpectCheck
	stateExpr
	stateExprIgnore
	stateEnd
)

var on = &sugar.LexemePredicate{}

func newRecognizer() sugar.StateMachine[state, node, node, node] {
	do := &nodeBuilder{}
	transitions := []sugar.Transition[state]{
		{From: stateStart, Event: on.Boundary, To: stateTarget, Action: do.clearLHS},
	}

	return &recognizer{transitions: transitions, builder: do}
}

type recognizer struct {
	transitions []sugar.Transition[state]
	builder     *nodeBuilder
}

func (r *recognizer) Transition(current state, lex sugar.Lexeme) (state, func(sugar.Lexeme)) {
	for _, row := range r.transitions {
		if row.From == current && row.Event(lex) {
			return row.To, row.Invoke()
		}
	}
	return stateStart, r.builder.reset
}

func (r *recognizer) InitialState() state {
	return stateStart
}

func (r *recognizer) Status(s state) sugar.Status {
	switch s {
	case stateIdle, stateStart, stateEnd:
		return sugar.StatusCompleted

	default:
		return sugar.StatusRunning
	}
}

func (r *recognizer) Build() node {
	return r.builder.build()
}

func (r *recognizer) BuildPartial() node {
	return r.builder.build()
}

func (r *recognizer) BuildError() node {
	return r.builder.build()
}

var _ sugar.StateMachine[state, node, node, node] = (*recognizer)(nil)
