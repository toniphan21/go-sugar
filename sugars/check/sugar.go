package check

import (
	"bytes"
	"strings"

	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/lex"
)

func New() sugar.Sugar {
	return &sugarImpl{Base: &sugar.Base{}}
}

const keyword = "check"
const ID = "nhatp.com/go/sugar/sugars/check"

type sugarImpl struct {
	*sugar.Base
}

func (s *sugarImpl) ID() string {
	return ID
}

func (s *sugarImpl) Parse(source []byte) []sugar.Node {
	return sugar.DoParse(LexicalParser(), source)
}

func (s *sugarImpl) Restore(source []byte) []sugar.Node {
	return sugar.DoParse(lex.SugarPlaceholderFuncParser(keyword), source)
}

func (s *sugarImpl) StructuralTransform(sourceID string, n sugar.Node) ([]byte, error) {
	return sugar.DoTransform[*node](s.Base, sourceID, n, func(source []byte, data *node) ([]byte, error) {
		// input:  [...]check ...
		// output: [...]__sugar_check__(...)
		// keep: pos-checkPos -> replace "check " by "__sugar_check__(" -> checkEnd:end -> add ")"
		out := bytes.Buffer{}

		out.Write(source[data.pos.Offset:data.checkPos.Offset])
		out.Write([]byte(lex.SugarPlaceholderFuncName(keyword)))
		out.WriteRune('(')
		out.Write(source[data.checkEnd.Offset:data.end.Offset])
		out.WriteRune(')')

		return out.Bytes(), nil
	})
}

func (s *sugarImpl) SemanticTransformer(sourceID string, scopeID string, n sugar.Node) ([]byte, error) {
	return sugar.DoTransformWithSemanticScope[*node](s.Base, sourceID, scopeID, n, func(source []byte, data *node) ([]byte, error) {
		// do simple things first, we update edge cases later
		// output:

		// if no identifiers:
		// if err := [checkEnd:end); err != nil {<nl>
		// <tab>return err<nl>
		// }<nl>

		// if there is identifiers:
		// [identifiers, ] err := [checkEnd:end)\n
		// if err != nil {<nl>
		// <tab>return err<nl>
		// }<nl>
		out := bytes.Buffer{}

		if len(data.identifiers) == 0 {
			out.WriteString("if err := ")
			out.Write(source[data.checkEnd.Offset:data.end.Offset])
			out.WriteString("; err != nil {\n")
			out.WriteRune('\t')
			out.WriteString("return err\n")
			out.WriteRune('}')
		} else {
			idents := make([]string, len(data.identifiers)+1)
			copy(idents, data.identifiers)
			idents[len(data.identifiers)] = "err"

			out.WriteString(strings.Join(idents, ", "))
			out.WriteString(" := ")
			out.Write(source[data.checkEnd.Offset:data.end.Offset])
			out.WriteRune('\n')
			out.WriteString("if err != nil {\n")
			out.WriteRune('\t')
			out.WriteString("return err\n")
			out.WriteRune('}')
		}

		return out.Bytes(), nil
	})
}

func (s *sugarImpl) RestoreTransform(sourceID string, n sugar.Node) ([]byte, error) {
	return sugar.DoTransform[lex.SugarPlaceholderFunc](s.Base, sourceID, n, func(source []byte, data lex.SugarPlaceholderFunc) ([]byte, error) {
		if data.Keyword() != keyword {
			return nil, sugar.ErrUnknownNode
		}

		// input:  [...]__sugar_check__(...)
		// output: [...]check ...
		// replace "__sugar_check__(" by "check ", keep: innerPos -> innerEnd,

		out := bytes.Buffer{}
		out.WriteString("check ")
		out.Write(source[data.InnerPos().Offset:data.InnerEnd().Offset])

		return out.Bytes(), nil
	})
}

var _ sugar.Sugar = (*sugarImpl)(nil)
