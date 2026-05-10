package sugar

var doNothing = func(Lexeme) {}

type transition[S comparable] interface {
	reset()

	invoke(current S, lexemes []Lexeme) (S, func(Lexeme), int, bool)
}

// ---

type TransitionTable[S comparable] interface {
	Route(from S, handler func(lex Lexeme) (S, bool)) TransitionTable[S]

	Add(from S, event func(Lexeme) bool, to S, actions ...func(lexeme Lexeme)) TransitionTable[S]

	Use(from S, parser LexicalParser, handle TransitionControl[S]) TransitionTable[S]

	Longest(from S, control TransitionControl[S], parsers ...LexicalParser) TransitionTable[S]

	Invoke(current S, lexemes []Lexeme) (S, func(lex Lexeme), int)
}

func NewTransitionTable[S comparable]() TransitionTable[S] {
	return &transitionTableImpl[S]{}
}

type transitionTableImpl[S comparable] struct {
	transitions []transition[S]
}

func (t *transitionTableImpl[S]) Route(from S, handler func(lex Lexeme) (S, bool)) TransitionTable[S] {
	t.transitions = append(t.transitions, &routeTransition[S]{
		from:    from,
		handler: handler,
	})
	return t
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

func (t *transitionTableImpl[S]) Longest(from S, control TransitionControl[S], parser ...LexicalParser) TransitionTable[S] {
	t.transitions = append(t.transitions, &longestWinTransition[S]{
		from:    from,
		parsers: parser,
		done:    make([]bool, len(parser)),
		control: control,
	})
	return t
}

func (t *transitionTableImpl[S]) Invoke(current S, lexemes []Lexeme) (S, func(lex Lexeme), int) {
	for _, v := range t.transitions {
		if next, action, consumed, ok := v.invoke(current, lexemes); ok {
			v.reset()

			return next, action, consumed
		}
	}
	return current, doNothing, 1
}

// ---

type routeTransition[S comparable] struct {
	from    S
	handler func(lex Lexeme) (S, bool)
}

func (t *routeTransition[S]) reset() {
}

func (t *routeTransition[S]) invoke(current S, lexemes []Lexeme) (S, func(lex Lexeme), int, bool) {
	if t.from == current && len(lexemes) != 0 && t.handler != nil {
		if dest, ok := t.handler(lexemes[0]); ok {
			return dest, doNothing, 1, true
		}
	}
	return current, doNothing, 0, false
}

// ---

type nodeTransition[S comparable] struct {
	from    S
	event   func(lex Lexeme) bool
	to      S
	actions []func(lex Lexeme)
}

func (t *nodeTransition[S]) reset() {
}

func (t *nodeTransition[S]) invoke(current S, lexemes []Lexeme) (S, func(lex Lexeme), int, bool) {
	if t.from == current && len(lexemes) != 0 && t.event(lexemes[0]) {
		return t.to, t.wrapAction(), 1, true
	}
	return current, doNothing, 0, false
}

func (t *nodeTransition[S]) wrapAction() func(Lexeme) {
	actions := t.actions
	return func(l Lexeme) {
		for _, a := range actions {
			a(l)
		}
	}
}

// ---

type TransitionControl[S comparable] struct {
	FirstTake func(lex Lexeme)

	WhenSuccess    func(p LexicalParser, d any, l Lexeme) (S, int)
	SuccessPutBack int
	SuccessMoveTo  S
	SuccessAction  func(data any, lex Lexeme)

	WhenError    func(p LexicalParser, d any, l Lexeme) (S, int)
	ErrorPutBack int
	ErrorMoveTo  S
	ErrorAction  func(data any, lex Lexeme)
}

func (c TransitionControl[S]) handleSuccess(parser LexicalParser, data any, lex Lexeme) (S, func(lex Lexeme), int) {
	if c.WhenSuccess != nil {
		next, putback := c.WhenSuccess(parser, data, lex)
		return next, doNothing, putback
	}

	if c.SuccessAction != nil {
		return c.SuccessMoveTo, func(lex Lexeme) { c.SuccessAction(data, lex) }, c.SuccessPutBack
	}
	return c.SuccessMoveTo, doNothing, c.SuccessPutBack
}

func (c TransitionControl[S]) handleError(parser LexicalParser, data any, lex Lexeme) (S, func(lex Lexeme), int) {
	if c.WhenError != nil {
		next, putback := c.WhenError(parser, data, lex)
		return next, doNothing, putback
	}

	if c.ErrorAction != nil {
		return c.ErrorMoveTo, func(lex Lexeme) { c.ErrorAction(data, lex) }, c.ErrorPutBack
	}
	return c.ErrorMoveTo, doNothing, c.ErrorPutBack
}

type lexicalParserTransition[S comparable] struct {
	from      S
	parser    LexicalParser
	control   TransitionControl[S]
	firstTook bool
}

func (t *lexicalParserTransition[S]) reset() {
	t.firstTook = false
	t.parser.Reset()
}

func (t *lexicalParserTransition[S]) invoke(current S, lexemes []Lexeme) (S, func(lex Lexeme), int, bool) {
	if t.from != current || len(lexemes) == 0 {
		return current, doNothing, 0, false
	}

	if !t.firstTook && t.control.FirstTake != nil {
		t.control.FirstTake(lexemes[0])
		t.firstTook = true
	}

	if t.parser.Done(lexemes[0:1]) { // this transition take slice with 1 lex at a time
		data, ok := t.parser.Result()

		if ok {
			next, action, putback := t.control.handleSuccess(t.parser, data, lexemes[0])
			return next, action, 1 - putback, true
		}

		next, action, putback := t.control.handleError(t.parser, data, lexemes[0])
		return next, action, 1 - putback, true
	}
	return current, doNothing, 1, false
}

// ---

type longestWinTransition[S comparable] struct {
	from      S
	parsers   []LexicalParser
	done      []bool
	control   TransitionControl[S]
	firstTook bool
}

func (t *longestWinTransition[S]) reset() {
	t.firstTook = false
	for i := range t.done {
		t.done[i] = false
	}
	for _, p := range t.parsers {
		p.Reset()
	}
}

func (t *longestWinTransition[S]) longestParser(lexemes []Lexeme) (LexicalParser, any, bool, int) {
	consumed := 0
	for i := range lexemes {
		allDone := true
		for j, p := range t.parsers {
			if t.done[j] {
				continue
			}
			if p.Done(lexemes[i : i+1]) {
				t.done[j] = true
			} else {
				allDone = false
			}
		}
		consumed++
		if allDone {
			break
		}
	}

	for _, d := range t.done {
		if !d {
			return nil, nil, false, consumed // not all done yet
		}
	}

	var bestData any
	bestIdx := -1
	bestConsumed := -1
	for i, p := range t.parsers {
		if data, ok := p.Result(); ok {
			if c := p.Consumed(); c > bestConsumed {
				bestConsumed = c
				bestIdx = i
				bestData = data
			}
		}
	}

	if bestIdx == -1 {
		return nil, nil, true, consumed // all done, no winner
	}
	return t.parsers[bestIdx], bestData, true, consumed // all done, there is a winner
}

func (t *longestWinTransition[S]) invoke(current S, lexemes []Lexeme) (S, func(lex Lexeme), int, bool) {
	if t.from != current || len(lexemes) == 0 {
		return current, doNothing, 0, false
	}

	if !t.firstTook && t.control.FirstTake != nil {
		t.control.FirstTake(lexemes[0])
		t.firstTook = true
	}

	parser, data, done, consumed := t.longestParser(lexemes)
	if done {
		if parser != nil {
			lastLex := lexemes[parser.Consumed()-1]
			next, action, putback := t.control.handleSuccess(parser, data, lastLex)
			return next, action, parser.Consumed() - putback, true
		}

		lastLex := lexemes[consumed-1]
		next, action, putback := t.control.handleError(nil, data, lastLex)
		return next, action, consumed - putback, true
	}
	return current, doNothing, consumed, false
}
