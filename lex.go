package sugar

import (
	"fmt"
	"go/scanner"
	"go/token"
	"log/slog"
)

type Lexeme struct {
	Tok    token.Token
	Pos    token.Pos
	Lit    string
	Line   int
	Column int
	Offset int
}

func (l *Lexeme) GoString() string {
	return fmt.Sprintf("%v | Tok:%d Pos:%d Lit:%#v | line:%d col:%d offset:%d", TokenName(l.Tok), l.Tok, l.Pos, l.Lit, l.Line, l.Column, l.Offset)
}

func Lex(content []byte) []Lexeme {
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(content))

	var s scanner.Scanner
	s.Init(file, content, nil, 0)

	var result []Lexeme
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}

		position := fset.Position(pos)
		lex := Lexeme{
			Tok:    tok,
			Pos:    pos,
			Lit:    lit,
			Line:   position.Line,
			Column: position.Column,
			Offset: position.Offset,
		}
		if lex.Lit == "" {
			lex.Lit = tok.String()
		}
		result = append(result, lex)
	}
	return result
}

// ---

type Status int

type LexicalNodeBuilder[N any] interface {
	SetParserID(name string)

	Reset()

	Build() (N, bool)
}

type LexicalParser interface {
	ID() string

	Is(parser LexicalParser) bool

	Reset()

	Done(lexemes []Lexeme) bool

	Result() (any, bool)

	Consumed() int
}

// ---

type lexicalParserImpl[S comparable, B LexicalNodeBuilder[N], N any] struct {
	id           string
	table        TransitionTable[S]
	current      S
	initialState S
	endState     S
	consumed     int
	builder      B
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
	i, total, lenLex := 0, 0, len(lexemes)
	for i < lenLex {
		idx := i
		from := p.current
		debug("invoke",
			slog.String("parser", p.id),
			slog.String("component", "LexicalParser"),
			slog.String("firstLex", lexemes[idx].GoString()),
			slog.Int("lexemesLength", lenLex),
			slog.Int("invokeLexemesIndex", idx),
			slog.Any("from", from),
			slog.Int("parserConsumed", p.consumed),
		)
		next, action, consumed := p.table.Invoke(from, lexemes[idx:])

		p.current = next
		action(lexemes[idx])

		total += consumed
		p.consumed += consumed
		i += consumed

		if p.current == p.endState {
			debug("done",
				slog.String("parser", p.id),
				slog.String("component", "LexicalParser"),
				slog.Int("lexemesLength", lenLex),
				slog.Int("invokeLexemesIndex", idx),
				slog.Int("invokeConsumed", consumed),
				slog.Any("from", from),
				slog.Any("to", p.endState),
				slog.Int("consumed", total),
				slog.Int("parserConsumed", p.consumed),
			)
			return true
		}

		debug("invoked",
			slog.String("parser", p.id),
			slog.String("component", "LexicalParser"),
			slog.Int("lexemesLength", lenLex),
			slog.Int("invokeLexemesIndex", idx),
			slog.Int("invokeConsumed", consumed),
			slog.Any("from", from),
			slog.Any("to", p.endState),
			slog.Int("consumed", total),
			slog.Int("parserConsumed", p.consumed),
		)
	}

	debug("not-done-yet",
		slog.String("parser", p.id),
		slog.String("component", "LexicalParser"),
		slog.Int("lexemesLength", lenLex),
		slog.Any("state", p.current),
		slog.Int("consumed", total),
		slog.Int("parserConsumed", p.consumed),
	)
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
	builder.SetParserID(id)

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

// ---

type Node interface {
	AsSugar() (Sugar, bool)
}

type NodeBuilder[T any] struct {
	Error    bool
	Node     *T
	onBuild  func(*T, bool)
	counters map[string]int
	parserID string
}

func NewNodeBuilder[T any]() *NodeBuilder[T] {
	return &NodeBuilder[T]{
		Node:     new(T),
		counters: make(map[string]int),
	}
}

func (b *NodeBuilder[T]) SetParserID(name string) {
	b.parserID = name
}

func (b *NodeBuilder[T]) SetName(name string) {
	b.parserID = name
}

func (b *NodeBuilder[T]) Reset() {
	b.Error = false
	b.Node = new(T)
	b.counters = make(map[string]int)
}

func (b *NodeBuilder[T]) Build() (T, bool) {
	ok := !b.Error
	if b.onBuild != nil {
		b.onBuild(b.Node, ok)
	}

	node := *b.Node
	msg := "build-success"
	if !ok {
		msg = "build-fail"
	}
	debug(msg,
		slog.String("parser", b.parserID),
		slog.Bool("ok", ok),
		slog.String("component", "NodeBuilder"),
		slog.Any("node", node),
	)
	return node, ok
}

func (b *NodeBuilder[T]) CounterInc(name string) {
	val, have := b.counters[name]
	if !have {
		b.counters[name] = 1
	} else {
		b.counters[name] = val + 1
	}
}

func (b *NodeBuilder[T]) CounterDec(name string) {
	val, have := b.counters[name]
	if !have {
		b.counters[name] = -1
	} else {
		b.counters[name] = val - 1
	}
}

func (b *NodeBuilder[T]) Counter(name string) int {
	val, have := b.counters[name]
	if !have {
		return 0
	}
	return val
}

func (b *NodeBuilder[T]) OnBuild(fn func(*T, bool)) *NodeBuilder[T] {
	b.onBuild = fn
	return b
}

func (b *NodeBuilder[T]) Fail(lex Lexeme) {
	b.Error = true
	debug("fail",
		slog.String("parser", b.parserID),
		slog.String("component", "NodeBuilder"),
		slog.Any("error", b.Error),
		slog.String("lexeme", lex.GoString()),
	)
}

func (b *NodeBuilder[T]) Collect(name string, fn func(n *T, l Lexeme)) func(Lexeme) {
	return func(lex Lexeme) {
		be := b.Error
		bn := *b.Node
		fn(b.Node, lex)
		debug("collect",
			slog.String("parser", b.parserID),
			slog.String("component", "NodeBuilder"),
			slog.Bool("beforeError", be),
			slog.Any("beforeNode", bn),
			slog.Bool("afterError", b.Error),
			slog.Any("afterNode", *b.Node),
		)
	}
}

func (b *NodeBuilder[T]) FailInner(inner any, lex Lexeme) {
	b.Error = true
	debug("fail-inner",
		slog.String("parser", b.parserID),
		slog.String("component", "NodeBuilder"),
		slog.Any("error", b.Error),
		slog.Any("data", inner),
		slog.String("lexeme", lex.GoString()),
	)
}

func (b *NodeBuilder[T]) CollectInner(name string, fn func(n *T, i any, l Lexeme)) func(any, Lexeme) {
	return func(data any, lex Lexeme) {
		be := b.Error
		bn := *b.Node
		fn(b.Node, data, lex)
		debug("collect-inner",
			slog.String("parser", b.parserID),
			slog.String("component", "NodeBuilder"),
			slog.String("name", name),
			slog.Any("data", data),
			slog.Bool("beforeError", be),
			slog.Any("beforeNode", bn),
			slog.Bool("afterError", b.Error),
			slog.Any("afterNode", *b.Node),
		)
	}
}
