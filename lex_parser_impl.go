//go:build !dev

package sugar

type lexicalParserImpl[S comparable, B LexicalNodeBuilder[N], N any] struct {
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

func NewLexicalParser[S comparable, B LexicalNodeBuilder[N], N any](
	name string,
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
