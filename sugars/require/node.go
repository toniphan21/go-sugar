package require

import (
	"golang.org/x/tools/go/packages"
	"nhatp.com/go/sugar"
)

type node struct {
	pos         sugar.Lexeme
	end         sugar.Lexeme
	requirePos  sugar.Lexeme
	requireEnd  sugar.Lexeme
	messagePos  *sugar.Lexeme
	messageEnd  *sugar.Lexeme
	message     *string
	identifiers []string
}

func (s *node) Pos() sugar.Lexeme {
	return s.pos
}

func (s *node) End() sugar.Lexeme {
	return s.end
}

func (s *node) SemanticAnalysis(pkg *packages.Package, smap *sugar.SourceMap) error {
	return nil
}

var _ sugar.SemanticNode = (*node)(nil)
