//go:build dev

package sugar

import "fmt"

type NodeBuilder[T any] struct {
	Error   bool
	Node    *T
	onBuild func(*T, bool)
	name    string
	debug   bool
}

func NewNodeBuilder[T any]() *NodeBuilder[T] {
	return &NodeBuilder[T]{
		Node: new(T),
	}
}

func (b *NodeBuilder[T]) SetName(name string) {
	b.name = name
}

func (b *NodeBuilder[T]) Debug() {
	b.debug = true
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
	if b.debug {
		fmt.Printf("%s:build: data=%#v ok=%v\n", b.name, node, ok)
	}
	return node, ok
}

func (b *NodeBuilder[T]) OnBuild(fn func(*T, bool)) *NodeBuilder[T] {
	b.onBuild = fn
	return b
}

func (b *NodeBuilder[T]) Fail(lex Lexeme) {
	b.Error = true
	if b.debug {
		fmt.Printf("%s:fail: error=%-5t lex=%#v\n", b.name, b.Error, &lex)
	}
}

func (b *NodeBuilder[T]) Collect(name string, fn func(n *T, l Lexeme)) func(Lexeme) {
	return func(lex Lexeme) {
		be := b.Error
		bn := *b.Node
		fn(b.Node, lex)
		ae := b.Error
		an := *b.Node

		if b.debug {
			fmt.Printf("%s:collect:%v lex=%#v\n\tbefore error=%-5t node=%#v\n\t after error=%-5t node=%#v\n", b.name, name, &lex, be, bn, ae, an)
		}
	}
}

func (b *NodeBuilder[T]) FailInner(inner any, lex Lexeme) {
	b.Error = true
	if b.debug {
		fmt.Printf("%s:fail: error=%-5t inner=%#v lex=%#v\n", b.name, b.Error, inner, &lex)
	}
}

func (b *NodeBuilder[T]) CollectInner(name string, fn func(n *T, i any, l Lexeme)) func(any, Lexeme) {
	return func(data any, lex Lexeme) {
		be := b.Error
		bn := *b.Node
		fn(b.Node, data, lex)
		ae := b.Error
		an := *b.Node
		if b.debug {
			fmt.Printf("%s:collect:%v data=%#v lex=%#v\n\tbefore error=%-5t node=%#v\n\t after error=%-5t node=%#v\n", b.name, name, data, &lex, be, bn, ae, an)
		}
	}
}
