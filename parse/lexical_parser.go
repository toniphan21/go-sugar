package parse

import "nhatp.com/go/sugar"

type LexicalParser[N any] interface {
	Reset()

	Take(lex sugar.Lexeme) sugar.Status

	Build() N
}

type builder[N any] interface {
	reset()

	build() N
}

type lexicalParserImpl[S ~int, B builder[N], N any] struct {
	transitions  []sugar.Transition[S]
	current      S
	builder      B
	initialState S
	statusMapper func(S, B) sugar.Status
}

func (p *lexicalParserImpl[S, B, N]) Reset() {
	p.current = p.initialState
	p.builder.reset()
}

func (p *lexicalParserImpl[S, B, N]) Take(lex sugar.Lexeme) sugar.Status {
	for _, row := range p.transitions {
		if row.From == p.current && row.Event(lex) {
			p.current = row.To
			row.Invoke()(lex)
			break
		}
	}
	return p.statusMapper(p.current, p.builder)
}

func (p *lexicalParserImpl[S, B, N]) Build() N {
	return p.builder.build()
}
