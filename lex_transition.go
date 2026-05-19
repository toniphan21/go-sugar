package sugar

import "log/slog"

var doNothing = func(Lexeme) {}

type transition[S comparable] interface {
	reset()

	invoke(current S, lexemes []Lexeme) (S, func(Lexeme), int, bool)
}

// ---

type transitionTableBuilder[S comparable] interface {
	// Route adds a transition decided by a callback handler.
	// The handler inspects the input and returns the next state. No putback is performed.
	Route(from S, handler func(lex Lexeme) (S, bool)) TransitionTable[S]

	// Scout adds a transition decided by a callback handler.
	// The handler inspects the input, returns the next state. The lexeme will be put back.
	Scout(from S, handler func(lex Lexeme) (S, bool)) TransitionTable[S]

	// Add adds a transition that checks the next lexeme against a predicate.
	// If matched, the parser advances to the given state without putback.
	Add(from S, event func(Lexeme) bool, to S, actions ...func(lexeme Lexeme)) TransitionTable[S]

	// Peek adds a transition that checks the next lexeme against a predicate.
	// If matched, the parser advances to the given state and puts back the lexeme.
	Peek(from S, event func(Lexeme) bool, to S, actions ...func(lexeme Lexeme)) TransitionTable[S]

	// Use delegates parsing to a sub LexicalParser with a transition control
	Use(from S, parser LexicalParser, control TransitionControl[S]) TransitionTable[S]

	// Optional delegates parsing to a sub LexicalParser.
	// On failure, all consumed input is put back and the transition is skipped.
	Optional(from S, parser LexicalParser, control OptionalTransitionControl[S]) TransitionTable[S]

	// Longest tries multiple LexicalParser and picks the one that consumes the most input.
	Longest(from S, control TransitionControl[S], parsers ...LexicalParser) TransitionTable[S]

	//// LongestOptional tries multiple sub LexicalParser and picks the longest match.
	//// If no parser matches, all consumed input is put back and the transition is skipped.
	//LongestOptional(from S, control OptionalTransitionControl[S], parsers ...LexicalParser) TransitionTable[S]
}

type TransitionTable[S comparable] interface {
	transitionTableBuilder[S]

	Invoke(current S, lexemes []Lexeme) (S, func(lex Lexeme), int)
}

func NewTransitionTable[S comparable](parserID string) TransitionTable[S] {
	return &transitionTableImpl[S]{
		parserID: parserID,
	}
}

type transitionTableImpl[S comparable] struct {
	parserID    string
	transitions []transition[S]
}

func (t *transitionTableImpl[S]) Route(from S, handler func(lex Lexeme) (S, bool)) TransitionTable[S] {
	t.transitions = append(t.transitions, &routeTransition[S]{
		parserID:   t.parserID,
		from:       from,
		handler:    handler,
		putback:    0,
		transition: "Route",
	})
	return t
}

func (t *transitionTableImpl[S]) Scout(from S, handler func(lex Lexeme) (S, bool)) TransitionTable[S] {
	t.transitions = append(t.transitions, &routeTransition[S]{
		parserID:   t.parserID,
		from:       from,
		handler:    handler,
		putback:    1,
		transition: "Scout",
	})
	return t
}

func (t *transitionTableImpl[S]) Add(from S, event func(Lexeme) bool, to S, actions ...func(lexeme Lexeme)) TransitionTable[S] {
	t.transitions = append(t.transitions, &nodeTransition[S]{
		parserID:   t.parserID,
		from:       from,
		event:      event,
		to:         to,
		actions:    actions,
		putback:    0,
		transition: "Add",
	})
	return t
}

func (t *transitionTableImpl[S]) Peek(from S, event func(Lexeme) bool, to S, actions ...func(lexeme Lexeme)) TransitionTable[S] {
	t.transitions = append(t.transitions, &nodeTransition[S]{
		parserID:   t.parserID,
		from:       from,
		event:      event,
		to:         to,
		actions:    actions,
		putback:    1,
		transition: "Peek",
	})
	return t
}

func (t *transitionTableImpl[S]) Use(from S, parser LexicalParser, control TransitionControl[S]) TransitionTable[S] {
	t.transitions = append(t.transitions, &lexicalParserTransition[S]{
		from:    from,
		parser:  parser,
		control: &control,
	})
	return t
}

func (t *transitionTableImpl[S]) Optional(from S, parser LexicalParser, control OptionalTransitionControl[S]) TransitionTable[S] {
	t.transitions = append(t.transitions, &lexicalParserTransition[S]{
		from:            from,
		parser:          parser,
		optionalControl: &control,
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
	parserID   string
	from       S
	handler    func(lex Lexeme) (S, bool)
	transition string
	putback    int
}

func (t *routeTransition[S]) reset() {
}

func (t *routeTransition[S]) invoke(current S, lexemes []Lexeme) (S, func(lex Lexeme), int, bool) {
	if t.from == current && len(lexemes) != 0 && t.handler != nil {
		if dest, ok := t.handler(lexemes[0]); ok {
			debug("transition-success",
				slog.String("parser", t.parserID),
				slog.String("component", "routeTransition"),
				slog.String("transition", t.transition),
				slog.Any("from", t.from),
				slog.Any("to", dest),
				slog.Int("putback", t.putback),
				slog.Int("consumed", 1-t.putback),
			)
			return dest, doNothing, 1 - t.putback, true
		}
	}
	return current, doNothing, 0, false
}

// ---

type nodeTransition[S comparable] struct {
	parserID   string
	from       S
	event      func(lex Lexeme) bool
	to         S
	actions    []func(lex Lexeme)
	transition string
	putback    int
}

func (t *nodeTransition[S]) reset() {
}

func (t *nodeTransition[S]) invoke(current S, lexemes []Lexeme) (S, func(lex Lexeme), int, bool) {
	if t.from == current && len(lexemes) != 0 && t.event(lexemes[0]) {
		debug("transition-success",
			slog.String("parser", t.parserID),
			slog.String("component", "nodeTransition"),
			slog.String("transition", t.transition),
			slog.Any("from", t.from),
			slog.Any("to", t.to),
			slog.Int("putback", t.putback),
			slog.Int("consumed", 1-t.putback),
		)
		return t.to, t.wrapAction(), 1 - t.putback, true
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

type OptionalTransitionControl[S comparable] struct {
	FirstTake     func(lex Lexeme)
	MoveTo        S
	SuccessAction func(data any, lex Lexeme)
	ErrorAction   func(data any, lex Lexeme)
}

func (c OptionalTransitionControl[S]) handleSuccess(data any) (S, func(lex Lexeme)) {
	if c.SuccessAction != nil {
		return c.MoveTo, func(lex Lexeme) { c.SuccessAction(data, lex) }
	}
	return c.MoveTo, doNothing
}

func (c OptionalTransitionControl[S]) handleError(data any) (S, func(lex Lexeme)) {
	if c.ErrorAction != nil {
		return c.MoveTo, func(lex Lexeme) { c.ErrorAction(data, lex) }
	}
	return c.MoveTo, doNothing
}

// ---

type lexicalParserTransition[S comparable] struct {
	from            S
	parser          LexicalParser
	control         *TransitionControl[S]
	optionalControl *OptionalTransitionControl[S]
	firstTook       bool
}

func (t *lexicalParserTransition[S]) reset() {
	t.firstTook = false
	t.parser.Reset()
}

func (t *lexicalParserTransition[S]) invoke(current S, lexemes []Lexeme) (S, func(lex Lexeme), int, bool) {
	if t.control != nil {
		return t.invokeRequired(current, lexemes)
	}
	return t.invokeOptional(current, lexemes)
}

func (t *lexicalParserTransition[S]) invokeRequired(current S, lexemes []Lexeme) (S, func(lex Lexeme), int, bool) {
	if t.from != current || len(lexemes) == 0 {
		return current, doNothing, 0, false
	}

	if !t.firstTook && t.control.FirstTake != nil {
		debug("first-take",
			slog.String("parser", t.parser.ID()),
			slog.String("component", "lexicalParserTransition"),
			slog.String("transition", "Use"),
		)
		t.control.FirstTake(lexemes[0])
		t.firstTook = true
	}

	if t.parser.Done(lexemes) {
		data, ok := t.parser.Result()
		consumed := t.parser.Consumed()
		lastLex := lexemes[consumed-1]

		if ok {
			next, action, putback := t.control.handleSuccess(t.parser, data, lastLex)
			debug("transition-success",
				slog.String("parser", t.parser.ID()),
				slog.String("component", "lexicalParserTransition"),
				slog.String("transition", "Use"),
				slog.Any("from", t.from),
				slog.Any("to", next),
				slog.Int("putback", putback),
				slog.Int("consumed", consumed-putback),
			)
			return next, action, consumed - putback, true
		}

		next, action, putback := t.control.handleError(t.parser, data, lastLex)
		debug("transition-fail",
			slog.String("parser", t.parser.ID()),
			slog.String("component", "lexicalParserTransition"),
			slog.String("transition", "Use"),
			slog.Any("from", t.from),
			slog.Any("to", next),
			slog.Int("putback", putback),
			slog.Int("consumed", consumed-putback),
		)
		return next, action, consumed - putback, true
	}
	return current, doNothing, 0, false
}

func (t *lexicalParserTransition[S]) invokeOptional(current S, lexemes []Lexeme) (S, func(lex Lexeme), int, bool) {
	if t.from != current || len(lexemes) == 0 {
		return current, doNothing, 0, false
	}

	if !t.firstTook && t.optionalControl.FirstTake != nil {
		debug("first-take",
			slog.String("parser", t.parser.ID()),
			slog.String("component", "lexicalParserTransition"),
			slog.String("transition", "Optional"),
		)
		t.optionalControl.FirstTake(lexemes[0])
		t.firstTook = true
	}

	if t.parser.Done(lexemes) {
		data, ok := t.parser.Result()

		if ok {
			next, action := t.optionalControl.handleSuccess(data)
			consumed := t.parser.Consumed()
			debug("transition-success",
				slog.String("parser", t.parser.ID()),
				slog.String("component", "lexicalParserTransition"),
				slog.String("transition", "Optional"),
				slog.Any("from", t.from),
				slog.Any("to", next),
				slog.Int("consumed", consumed),
			)
			return next, action, consumed, true
		}

		next, action := t.optionalControl.handleError(data)
		debug("transition-fail",
			slog.String("parser", t.parser.ID()),
			slog.String("component", "lexicalParserTransition"),
			slog.String("transition", "Optional"),
			slog.Any("from", t.from),
			slog.Any("to", next),
			slog.Int("consumed", 0),
		)
		return next, action, 0, true // if error just handle off and mark the consumed token as 0
	}
	return current, doNothing, 0, false
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
		if parser != nil { // winner
			lastLex := lexemes[parser.Consumed()-1]
			next, action, putback := t.control.handleSuccess(parser, data, lastLex)
			return next, action, parser.Consumed() - putback, true
		}

		// no-winner: error
		lastLex := lexemes[consumed-1]
		next, action, putback := t.control.handleError(nil, data, lastLex)
		return next, action, consumed - putback, true
	}
	return current, doNothing, consumed, false
}
