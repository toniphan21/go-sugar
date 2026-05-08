package sugar

type Status int

type LexicalNodeBuilder[N any] interface {
	Reset()

	Build() (N, bool)
}

type LexicalParser interface {
	Debug() LexicalParser

	ID() string

	Is(parser LexicalParser) bool

	Reset()

	Done(lexemes []Lexeme) bool

	Result() (any, bool)

	Consumed() int
}
