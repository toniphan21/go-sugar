package sugar

import (
	"bytes"
	"cmp"
	"slices"

	"golang.org/x/tools/go/packages"
)

type Sugar interface {
	Pos() Lexeme
	End() Lexeme

	StructuralTransform(source []byte, lexemes []Lexeme) []byte

	SemanticAnalysis(pkg *packages.Package, smap *SourceMap) error

	SemanticTransformer(source []byte, lexemes []Lexeme) []byte
}

type Plugin interface {
	ID() string

	LexicalParser() LexicalParser

	RestoreParser() LexicalParser

	RestoreTransform(node Node, source []byte, lexemes []Lexeme) []byte
}

var plugins map[string]Plugin

func Register(plugin Plugin) {
	if plugin == nil {
		return
	}
	if plugins == nil {
		plugins = make(map[string]Plugin)
	}
	plugins[plugin.ID()] = plugin
}

func asSugar(v any) (Sugar, bool) {
	if s, ok := v.(Sugar); ok {
		return s, true
	}

	if n, ok := v.(ParsedNode); ok {
		return n.AsSugar()
	}
	return nil, false
}

func parseSugars(lexemes []Lexeme) []Sugar {
	var sugars []Sugar
	for _, plugin := range plugins {
		parser := plugin.LexicalParser()

		offset := 0
		for offset < len(lexemes) {
			slice := lexemes[offset:]

			if parser.Done(slice) {
				if result, success := parser.Result(); success {
					if sugar, ok := asSugar(result); ok {
						sugars = append(sugars, sugar)
					}
				}
			}
			offset += parser.Consumed()
			parser.Reset()
		}
	}
	return sugars
}

type restoreNode struct {
	plugin Plugin
	node   Node
}

func parseRestoreNodes(lexemes []Lexeme) []*restoreNode {
	var nodes []*restoreNode
	var parsers = make(map[string]LexicalParser)
	for _, plugin := range plugins {
		parser, have := parsers[plugin.ID()]
		if !have {
			parser = plugin.RestoreParser()
			parsers[plugin.ID()] = parser
		}

		offset := 0
		for offset < len(lexemes) {
			slice := lexemes[offset:]

			if parser.Done(slice) {
				if result, success := parser.Result(); success {
					if n, ok := result.(Node); ok {
						nodes = append(nodes, &restoreNode{plugin: plugin, node: n})
					}
				}
			}
			offset += parser.Consumed()
			parser.Reset()
		}
	}
	return nodes
}

func RestoreTransform(source []byte) []byte {
	lexemes := Lex(source)
	nodes := parseRestoreNodes(lexemes)

	slices.SortFunc(nodes, func(a, b *restoreNode) int {
		return cmp.Compare(a.node.Pos().Offset, b.node.Pos().Offset)
	})

	for _, v := range nodes {
		v.plugin.RestoreTransform(v.node, source, lexemes)
	}

	out := bytes.Buffer{}
	cursor := 0
	for _, v := range nodes {
		out.Write(source[cursor:v.node.Pos().Offset])
		transformed := v.plugin.RestoreTransform(v.node, source, lexemes)
		out.Write(transformed)
		cursor = v.node.End().Offset
	}
	out.Write(source[cursor:])

	return out.Bytes()
}
