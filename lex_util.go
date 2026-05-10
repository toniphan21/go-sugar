package sugar

import "go/token"

func CollectBuilderDataOrFail[N, D any](builder *NodeBuilder[N], data any, lex Lexeme, fn func(D)) {
	if v, ok := data.(D); ok {
		fn(v)
	} else {
		builder.Fail(lex)
	}
}

type LexemePredicate struct {
}

func (*LexemePredicate) StatementBoundary(lex Lexeme) bool {
	return lex.Tok == token.SEMICOLON || lex.Tok == token.LBRACE
}

func (*LexemePredicate) IdentMatch(lit string) func(Lexeme) bool {
	return func(lex Lexeme) bool {
		return lex.Tok == token.IDENT && lex.Lit == lit
	}
}

func (*LexemePredicate) Ident(lex Lexeme) bool {
	return lex.Tok == token.IDENT
}

func (*LexemePredicate) IsNotIdent(lex Lexeme) bool {
	return lex.Tok != token.IDENT
}

func (p *LexemePredicate) Comma(lex Lexeme) bool {
	return lex.Tok == token.COMMA
}

func (p *LexemePredicate) Period(lex Lexeme) bool {
	return lex.Tok == token.PERIOD
}

func (p *LexemePredicate) Any(lex Lexeme) bool {
	return true
}

func (p *LexemePredicate) Assign(lex Lexeme) bool {
	return lex.Tok == token.ASSIGN
}

func (p *LexemePredicate) Define(lex Lexeme) bool {
	return lex.Tok == token.DEFINE
}

func (p *LexemePredicate) LeftParen(lex Lexeme) bool {
	return lex.Tok == token.LPAREN
}

func (p *LexemePredicate) RightParen(lex Lexeme) bool {
	return lex.Tok == token.RPAREN
}
