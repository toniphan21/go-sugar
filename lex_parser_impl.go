//go:build !dev

package sugar

type lexicalParserImpl[S comparable, B LexicalNodeBuilder[N], N any] struct {
	id           string
	table        TransitionTable[S]
	current      S
	initialState S
	endState     S
	consumed     int
	builder      B
}

func (p *lexicalParserImpl[S, B, N]) Debug() LexicalParser {
	return p
}

func (p *lexicalParserImpl[S, B, N]) ID() string {
	return p.id
}

func (p *lexicalParserImpl[S, B, N]) Is(parser LexicalParser) bool {
	return p.ID() == parser.ID()
}

func (p *lexicalParserImpl[S, B, N]) Reset() {
	p.current = p.initialState
	p.consumed = 0
	p.builder.Reset()
}

func (p *lexicalParserImpl[S, B, N]) Done(lexemes []Lexeme) bool {
	i := 0
	for i < len(lexemes) {
		next, action, consumed := p.table.Invoke(p.current, lexemes[i:])

		p.current = next
		action(lexemes[i])

		p.consumed += consumed
		i += consumed

		if p.current == p.endState {
			return true
		}
	}
	return false
}

func (p *lexicalParserImpl[S, B, N]) Result() (any, bool) {
	return p.builder.Build()
}

func (p *lexicalParserImpl[S, B, N]) Consumed() int {
	return p.consumed
}

func NewLexicalParser[S comparable, B LexicalNodeBuilder[N], N any](
	id string,
	transitionTable TransitionTable[S],
	initialState S,
	endState S,
	builder B,
) LexicalParser {
	return &lexicalParserImpl[S, B, N]{
		id:           id,
		table:        transitionTable,
		current:      initialState,
		initialState: initialState,
		endState:     endState,
		consumed:     0,
		builder:      builder,
	}
}
