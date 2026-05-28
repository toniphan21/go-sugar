package check

import (
	"bytes"

	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/lex"
)

func New() sugar.Plugin {
	return &pluginImpl{}
}

const keyword = "check"
const PluginID = "nhatp.com/go/sugar/sugars/check"

type pluginImpl struct {
}

func (i *pluginImpl) ID() string {
	return PluginID
}

func (i *pluginImpl) LexicalParser() sugar.LexicalParser {
	return LexicalParser()
}

func (i *pluginImpl) RestoreParser() sugar.LexicalParser {
	return lex.SugarPlaceholderFuncParser(keyword)
}

func (i *pluginImpl) RestoreTransform(node sugar.Node, source []byte, lexemes []sugar.Lexeme) []byte {
	// input:  [...]__sugar_check__(...)
	// output: [...]check ...
	// replace "__sugar_check__(" by "check ", keep: innerPos -> innerEnd,
	data, ok := node.(lex.SugarPlaceholderFunc)
	if !ok || data.Keyword() != keyword {
		return source[node.Pos().Offset:node.End().Offset]
	}

	out := bytes.Buffer{}
	out.WriteString("check ")
	out.Write(source[data.InnerPos().Offset:data.InnerEnd().Offset])

	return out.Bytes()
}

var _ sugar.Plugin = (*pluginImpl)(nil)
