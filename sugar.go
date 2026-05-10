package sugar

import (
	"bytes"
	"cmp"
	"slices"
)

type Sugar interface {
	Pos() Lexeme
	End() Lexeme

	StructuralTransform(source []byte, lexemes []Lexeme) []byte
}

type Plugin interface {
	ID() string

	LexicalParser() LexicalParser
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

// Execute is temporary name which triggers whole pipeline, will rename later
func Execute(source []byte) ([]byte, *SourceMap) {
	lexemes := Lex(source)
	sugars := doParseSugars(lexemes)

	slices.SortFunc(sugars, func(a, b Sugar) int {
		return cmp.Compare(a.Pos().Offset, b.Pos().Offset)
	})

	output, smap := doStructuralTransform(source, lexemes, sugars)
	return output, smap
}

func asSugar(v any) (Sugar, bool) {
	if s, ok := v.(Sugar); ok {
		return s, true
	}

	if n, ok := v.(Node); ok {
		return n.AsSugar()
	}
	return nil, false
}

func doParseSugars(lexemes []Lexeme) []Sugar {
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

func doStructuralTransform(source []byte, lexemes []Lexeme, sugars []Sugar) ([]byte, *SourceMap) {
	out := bytes.Buffer{}
	smap := &SourceMap{}
	cursor := 0

	for _, v := range sugars {
		out.Write(source[cursor:v.Pos().Offset])

		goStart := out.Len()
		transformed := v.StructuralTransform(source, lexemes)
		out.Write(transformed)
		goEnd := out.Len()

		smap.Entries = append(smap.Entries, Entry{
			Sugar: Region{
				Pos: Position{Offset: v.Pos().Offset},
				End: Position{Offset: v.End().Offset},
			},
			Go: Region{
				Pos: Position{Offset: goStart},
				End: Position{Offset: goEnd},
			},
			Kind: KindExpand,
		})

		cursor = v.End().Offset
	}

	out.Write(source[cursor:])

	return out.Bytes(), smap
}
