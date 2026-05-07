package sugar

var doNothing = func(Lexeme) {}

type transition[S comparable] interface {
	invoke(current S, lex Lexeme) (S, func(Lexeme), bool)
}

// ---

type TransitionTable[S comparable] interface {
	Add(from S, event func(Lexeme) bool, to S, actions ...func(lexeme Lexeme)) TransitionTable[S]

	Use(from S, parser LexicalParser, handle TransitionControl[S]) TransitionTable[S]

	Invoke(current S, lex Lexeme) (S, func(lex Lexeme))
}

func NewTransitionTable[S comparable]() TransitionTable[S] {
	return &transitionTableImpl[S]{}
}

type transitionTableImpl[S comparable] struct {
	transitions []transition[S]
}

func (t *transitionTableImpl[S]) Add(from S, event func(Lexeme) bool, to S, actions ...func(lexeme Lexeme)) TransitionTable[S] {
	t.transitions = append(t.transitions, &nodeTransition[S]{
		from:    from,
		event:   event,
		to:      to,
		actions: actions,
	})
	return t
}

func (t *transitionTableImpl[S]) Use(from S, parser LexicalParser, control TransitionControl[S]) TransitionTable[S] {
	t.transitions = append(t.transitions, &lexicalParserTransition[S]{
		from:    from,
		parser:  parser,
		control: control,
	})
	return t
}

func (t *transitionTableImpl[S]) Invoke(current S, lex Lexeme) (S, func(lex Lexeme)) {
	for _, v := range t.transitions {
		if next, action, ok := v.invoke(current, lex); ok {
			return next, action
		}
	}
	return current, doNothing
}

// ---

type nodeTransition[S comparable] struct {
	from    S
	event   func(lex Lexeme) bool
	to      S
	actions []func(lex Lexeme)
}

func (t *nodeTransition[S]) wrapAction() func(Lexeme) {
	actions := t.actions
	return func(l Lexeme) {
		for _, a := range actions {
			a(l)
		}
	}
}

func (t *nodeTransition[S]) invoke(current S, lex Lexeme) (S, func(lex Lexeme), bool) {
	if t.from == current && t.event(lex) {
		return t.to, t.wrapAction(), true
	}
	return current, func(Lexeme) {}, false
}

// ---

type TransitionControl[S comparable] struct {
	FirstTake     func(lex Lexeme)
	SuccessMoveTo S
	SuccessAction func(data any, lex Lexeme)
	ErrorMoveTo   S
	ErrorAction   func(data any, lex Lexeme)
}

type lexicalParserTransition[S comparable] struct {
	from     S
	parser   LexicalParser
	control  TransitionControl[S]
	consumed int
}

func (t *lexicalParserTransition[S]) invoke(current S, lex Lexeme) (S, func(lex Lexeme), bool) {
	if t.from != current {
		return current, doNothing, false
	}

	if t.consumed == 0 && t.control.FirstTake != nil {
		t.control.FirstTake(lex)
	}

	if t.parser.Done(lex) {
		data, ok := t.parser.Result()
		defer func() {
			t.consumed = 0
			t.parser.Reset()
		}()

		if ok {
			if t.control.SuccessAction != nil {
				return t.control.SuccessMoveTo, func(lex Lexeme) { t.control.SuccessAction(data, lex) }, true
			}
			return t.control.SuccessMoveTo, doNothing, true
		}

		if t.control.ErrorAction != nil {
			return t.control.ErrorMoveTo, func(lex Lexeme) { t.control.ErrorAction(data, lex) }, true
		}
		return t.control.ErrorMoveTo, doNothing, true
	}
	return current, doNothing, false
}
