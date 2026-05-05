package sugar

type Status int

const (
	StatusCompleted Status = iota
	StatusRunning
	StatusFailed
)

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
		case StatusCompleted:
			nodes = append(nodes, machine.Build())
			current = machine.InitialState()
		case StatusFailed:
			errors = append(errors, machine.BuildError())
			current = machine.InitialState()
		case StatusRunning:
			// keep going
		}
	}

	var partial P
	if machine.Status(current) == StatusRunning {
		partial = machine.BuildPartial()
	}
	return nodes, errors, partial
}

type LexicalNodeBuilder[N any] interface {
	Reset()

	Build() (N, bool)
}

type LexicalParser interface {
	Reset()

	Done(lex Lexeme) bool

	Result() (any, bool)

	Consumed() int
}

type lexicalParserImpl[S ~int, B LexicalNodeBuilder[N], N any] struct {
	table        TransitionTable[S]
	current      S
	initialState S
	endState     S
	consumed     int
	builder      B
}

func (p *lexicalParserImpl[S, B, N]) Reset() {
	p.current = p.initialState
	p.consumed = 0
	p.builder.Reset()
}

func (p *lexicalParserImpl[S, B, N]) Done(lex Lexeme) bool {
	next, action := p.table.Invoke(p.current, lex)

	p.consumed++
	p.current = next
	action(lex)

	return p.current == p.endState
}

func (p *lexicalParserImpl[S, B, N]) Result() (any, bool) {
	return p.builder.Build()
}

func (p *lexicalParserImpl[S, B, N]) Consumed() int {
	return p.consumed
}

func NewLexicalParser[S ~int, B LexicalNodeBuilder[N], N any](
	transitionTable TransitionTable[S],
	initialState S,
	endState S,
	builder B,
) LexicalParser {
	return &lexicalParserImpl[S, B, N]{
		table:        transitionTable,
		current:      initialState,
		initialState: initialState,
		endState:     endState,
		consumed:     0,
		builder:      builder,
	}
}
