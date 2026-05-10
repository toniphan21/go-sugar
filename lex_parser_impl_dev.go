//go:build dev

package sugar

import "fmt"

type NodeBuilderDev interface {
	SetName(string)

	Debug()
}

type lexicalParserImpl[S comparable, B LexicalNodeBuilder[N], N any] struct {
	id           string
	table        TransitionTable[S]
	current      S
	initialState S
	endState     S
	consumed     int
	builder      B
	debug        bool
}

func (p *lexicalParserImpl[S, B, N]) Debug() LexicalParser {
	if d, ok := any(p.builder).(NodeBuilderDev); ok {
		d.Debug()
	}
	p.debug = true

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

		if p.debug {
			fmt.Printf("%v[%v] consumed %d\n", p.id, p.current, p.consumed)
		}
		i += consumed

		if p.current == p.endState {
			if p.debug {
				fmt.Printf("%v[%v]->[%v]\n", p.id, p.current, p.endState)
			}
			return true
		}
	}

	if p.debug {
		fmt.Printf("%v[%v]: not done yet\n", p.id, p.current)
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
	if d, ok := any(builder).(NodeBuilderDev); ok {
		d.SetName(id)
	}

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
