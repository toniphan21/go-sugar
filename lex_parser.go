package sugar

type Status int

type LexicalNodeBuilder[N any] interface {
	Reset()

	Build() (N, bool)
}

type LexicalParser interface {
	Debug() LexicalParser

	Reset()

	Done(lex Lexeme) bool

	Result() (any, bool)

	Consumed() int
}
