package require

import (
	"bytes"

	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/lex"
)

func New() sugar.Plugin {
	return &pluginImpl{}
}

const keyword = "require"
const PluginID = "nhatp.com/go/sugar/sugars/require"

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
	data, ok := node.(lex.SugarPlaceholderFunc)
	if !ok || data.Keyword() != keyword {
		return source[node.Pos().Offset:node.End().Offset]
	}

	body := data.Body()
	lb := len(body)
	if lb == 0 {
		return source[node.Pos().Offset:node.End().Offset]
	}

	is := &sugar.LexemePredicate{}
	if lb > 2 && is.String(body[lb-1]) && is.Comma(body[lb-2]) {
		// end with message, we need to strip the comma before that
		out := bytes.Buffer{}
		out.WriteString("require ")
		out.Write(source[data.InnerPos().Offset:body[lb-2].Offset]) // inner pos til the comma
		out.WriteRune(' ')
		out.Write(source[body[lb-1].Offset:data.InnerEnd().Offset]) // literal string after comma

		return out.Bytes()
	}

	// no message, just unwrap the function
	out := bytes.Buffer{}
	out.WriteString("require ")
	out.Write(source[data.InnerPos().Offset:data.InnerEnd().Offset])

	return out.Bytes()
}

var _ sugar.Plugin = (*pluginImpl)(nil)
