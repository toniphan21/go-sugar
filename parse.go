package sugar

import "go/token"

type Status int

const (
	Terminal Status = iota
	Running
	Failed
)

type Transition[S ~int] struct {
	From    S
	Event   func(lex Lexeme) bool
	To      S
	Action  func(Lexeme)
	Actions []func(Lexeme)
}

func (t *Transition[S]) Invoke() func(Lexeme) {
	if t.Action != nil && t.Actions != nil {
		panic("transition: both Action and Actions set")
	}
	if t.Action != nil {
		return t.Action
	}

	actions := t.Actions
	return func(l Lexeme) {
		for _, a := range actions {
			a(l)
		}
	}
}

type StateMachine[S ~int, N, E, P any] interface {
	Transition(current S, lex Lexeme) (S, func(Lexeme))
	InitialState() S
	Status(S) Status
	Build() N
	BuildPartial() P
	BuildError() E
}

func RunStateMachine[S ~int, N, E, P any](machine StateMachine[S, N, E, P], lexemes []Lexeme) ([]N, []E, P) {
	current := machine.InitialState()
	var nodes []N
	var errors []E

	for _, lex := range lexemes {
		next, action := machine.Transition(current, lex)
		if action != nil {
			action(lex)
		}
		current = next

		switch machine.Status(current) {
		case Terminal:
			nodes = append(nodes, machine.Build())
			current = machine.InitialState()
		case Failed:
			errors = append(errors, machine.BuildError())
			current = machine.InitialState()
		case Running:
			// keep going
		}
	}

	var partial P
	if machine.Status(current) == Running {
		partial = machine.BuildPartial()
	}
	return nodes, errors, partial
}

type LexemePredicate struct {
}

func (*LexemePredicate) Boundary(lex Lexeme) bool {
	return lex.Tok == token.SEMICOLON
}

func (*LexemePredicate) Ident(lit string) func(Lexeme) bool {
	return func(lex Lexeme) bool {
		return lex.Tok == token.IDENT && lex.Lit == lit
	}
}
