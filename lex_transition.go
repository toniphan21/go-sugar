package sugar

type transition[S ~int] struct {
	from    S
	event   func(lex Lexeme) bool
	to      S
	actions []func(lex Lexeme)
}

func (t *transition[S]) wrapAction() func(Lexeme) {
	actions := t.actions
	return func(l Lexeme) {
		for _, a := range actions {
			a(l)
		}
	}
}

func (t *transition[S]) invoke(current S, lex Lexeme) (S, func(lex Lexeme), bool) {
	if t.from == current && t.event(lex) {
		return t.to, t.wrapAction(), true
	}
	return current, nil, false
}

// ---

type TransitionTable[S ~int] interface {
	Add(from S, event func(Lexeme) bool, to S, actions ...func(lexeme Lexeme)) TransitionTable[S]

	Invoke(current S, lex Lexeme) (S, func(lex Lexeme))
}

func NewTransitionTable[S ~int]() TransitionTable[S] {
	return &transitionTableImpl[S]{}
}

type transitionTableImpl[S ~int] struct {
	transitions []transition[S]
}

func (t *transitionTableImpl[S]) Add(from S, event func(Lexeme) bool, to S, actions ...func(lexeme Lexeme)) TransitionTable[S] {
	t.transitions = append(t.transitions, transition[S]{
		from:    from,
		event:   event,
		to:      to,
		actions: actions,
	})
	return t
}

func (t *transitionTableImpl[S]) Invoke(current S, lex Lexeme) (S, func(lex Lexeme)) {
	for _, v := range t.transitions {
		if next, action, ok := v.invoke(current, lex); ok {
			return next, action
		}
	}
	return current, func(lex Lexeme) {}
}

// ---
