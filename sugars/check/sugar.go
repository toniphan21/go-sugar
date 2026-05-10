package check

import (
	"bytes"

	"nhatp.com/go/sugar"
)

const StructuralReplacedName = "__sugar_check__("

type sugarImpl struct {
	pos         sugar.Lexeme
	end         sugar.Lexeme
	checkPos    sugar.Lexeme
	checkEnd    sugar.Lexeme
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
	// input:  [...]check ...
	// output: [...]__sugar_check__(...)
	// keep: pos-checkPos -> replace "check " by "__sugar_check__(" -> checkEnd:end -> add ")"
	out := bytes.Buffer{}

	out.Write(source[s.pos.Offset:s.checkPos.Offset])
	out.Write([]byte(StructuralReplacedName))
	out.Write(source[s.checkEnd.Offset:s.end.Offset])
	out.WriteRune(')')

	return out.Bytes()
}
