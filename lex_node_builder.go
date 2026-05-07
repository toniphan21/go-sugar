//go:build !dev

package sugar

type NodeBuilder[T any] struct {
	Error   bool
	Node    *T
	onBuild func(*T, bool)
}

func NewNodeBuilder[T any]() *NodeBuilder[T] {
	return &NodeBuilder[T]{
		Node: new(T),
	}
}

func (b *NodeBuilder[T]) Reset() {
	b.Error = false
	b.Node = new(T)
}

func (b *NodeBuilder[T]) Build() (T, bool) {
	ok := !b.Error
	if b.onBuild != nil {
		b.onBuild(b.Node, ok)
	}

	node := *b.Node
	return node, ok
}

func (b *NodeBuilder[T]) OnBuild(fn func(*T, bool)) *NodeBuilder[T] {
	b.onBuild = fn
	return b
}

func (b *NodeBuilder[T]) Fail(lex Lexeme) {
	b.Error = true
}

func (b *NodeBuilder[T]) Collect(name string, fn func(*T, Lexeme)) func(Lexeme) {
	return func(lex Lexeme) {
		fn(b.Node, lex)
	}
}

func (b *NodeBuilder[T]) FailInner(inner any, lex Lexeme) {}

func (b *NodeBuilder[T]) CollectInner(name string, fn func(n *T, i any, l Lexeme)) func(any, Lexeme) {
	return func(data any, lex Lexeme) {
		fn(b.Node, data, lex)
	}
}
