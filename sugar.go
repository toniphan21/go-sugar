package sugar

import "golang.org/x/tools/go/packages"

type Sugar interface {
	Pos() Lexeme
	End() Lexeme

	StructuralTransform(source []byte, lexemes []Lexeme) []byte

	SemanticAnalysis(pkg *packages.Package, smap *SourceMap) error
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

func asSugar(v any) (Sugar, bool) {
	if s, ok := v.(Sugar); ok {
		return s, true
	}

	if n, ok := v.(Node); ok {
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
