//go:build !dev

package sugar

type NodeBuilder[T any] struct {
	Error    bool
	Node     *T
	onBuild  func(*T, bool)
	counters map[string]int
	name     string
}

func NewNodeBuilder[T any]() *NodeBuilder[T] {
	return &NodeBuilder[T]{
		Node:     new(T),
		counters: make(map[string]int),
	}
}

func (b *NodeBuilder[T]) SetName(name string) {
	b.name = name
}

func (b *NodeBuilder[T]) Debug() {
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
}

func (b *NodeBuilder[T]) Collect(name string, fn func(n *T, l Lexeme)) func(Lexeme) {
	return func(lex Lexeme) {
		fn(b.Node, lex)
	}
}

func (b *NodeBuilder[T]) FailInner(inner any, lex Lexeme) {
	b.Error = true
}

func (b *NodeBuilder[T]) CollectInner(name string, fn func(n *T, i any, l Lexeme)) func(any, Lexeme) {
	return func(data any, lex Lexeme) {
		fn(b.Node, data, lex)
	}
}
