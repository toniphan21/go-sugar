package require

import (
	"bytes"

	"golang.org/x/tools/go/packages"
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/lex"
)

type sugarImpl struct {
	pos         sugar.Lexeme
	end         sugar.Lexeme
	requirePos  sugar.Lexeme
	requireEnd  sugar.Lexeme
	messagePos  *sugar.Lexeme
	messageEnd  *sugar.Lexeme
	message     *string
	identifiers []string
}

var _ sugar.Sugar = (*sugarImpl)(nil)

func (s *sugarImpl) Pos() sugar.Lexeme {
	return s.pos
}

func (s *sugarImpl) End() sugar.Lexeme {
	return s.end
}

func (s *sugarImpl) StructuralTransform(source []byte, _ []sugar.Lexeme) []byte {
	out := bytes.Buffer{}

	out.Write(source[s.pos.Offset:s.requirePos.Offset])
	out.Write([]byte(lex.SugarPlaceholderFuncName(keyword)))
	out.WriteRune('(')
	if s.message == nil {
		out.Write(source[s.requireEnd.Offset:s.end.Offset])
	} else {
		out.Write(source[s.requireEnd.Offset:s.messagePos.Offset])
		out.WriteRune(',')
		out.Write(source[s.messagePos.Offset:s.end.Offset])
	}
	out.WriteRune(')')

	return out.Bytes()
}

func (s *sugarImpl) SemanticTransformer(source []byte, lexemes []sugar.Lexeme) []byte {
	return source
}

func (s *sugarImpl) SemanticAnalysis(pkg *packages.Package, smap *sugar.SourceMap) error {
	return nil
}
