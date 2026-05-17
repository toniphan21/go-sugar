package sugar

import "go/token"

func TokenName(tok token.Token) string {
	t := int(tok)
	if t < len(tokenNames) {
		return tokenNames[t]
	} else {
		return tok.String()
	}
}

// copied and processed from token const in go/token package
var tokenNames = [...]string{
	// Special tokens
	"ILLEGAL",
	"EOF",
	"COMMENT",

	"literal_beg",
	// Identifiers and basic type literals
	// (these tokens stand for classes of literals)
	"IDENT",  // main
	"INT",    // "12345"",",
	"FLOAT",  // 123.45
	"IMAG",   // 123.45i
	"CHAR",   // 'a'
	"STRING", // "abc"
	"literal_end",

	"operator_beg",
	// Operators and delimiters
	"ADD", // +
	"SUB", // -
	"MUL", // *
	"QUO", // /
	"REM", // %

	"AND",     // &
	"OR",      // |
	"XOR",     // ^
	"SHL",     // <<
	"SHR",     // >>
	"AND_NOT", // &^

	"ADD_ASSIGN", // +=
	"SUB_ASSIGN", // -=
	"MUL_ASSIGN", // *=
	"QUO_ASSIGN", // /=
	"REM_ASSIGN", // %=

	"AND_ASSIGN",     // &=
	"OR_ASSIGN",      // |=
	"XOR_ASSIGN",     // ^=
	"SHL_ASSIGN",     // <<=
	"SHR_ASSIGN",     // >>=
	"AND_NOT_ASSIGN", // &^=

	"LAND",  // &&
	"LOR",   // ||
	"ARROW", // <-
	"INC",   // ++
	"DEC",   // --

	"EQL",    // ==
	"LSS",    // <
	"GTR",    // >
	"ASSIGN", // =
	"NOT",    // !

	"NEQ",      // !=
	"LEQ",      // <=
	"GEQ",      // >=
	"DEFINE",   // :=
	"ELLIPSIS", // ...

	"LPAREN", // (
	"LBRACK", // [
	"LBRACE", // {
	"COMMA",  // ,
	"PERIOD", // .

	"RPAREN",    // )
	"RBRACK",    // ]
	"RBRACE",    // }
	"SEMICOLON", // ;
	"COLON",     // :
	"operator_end",

	"keyword_beg",
	// Keywords
	"BREAK",
	"CASE",
	"CHAN",
	"CONST",
	"CONTINUE",

	"DEFAULT",
	"DEFER",
	"ELSE",
	"FALLTHROUGH",
	"FOR",

	"FUNC",
	"GO",
	"GOTO",
	"IF",
	"IMPORT",

	"INTERFACE",
	"MAP",
	"PACKAGE",
	"RANGE",
	"RETURN",

	"SELECT",
	"STRUCT",
	"SWITCH",
	"TYPE",
	"VAR",
	"keyword_end",

	"additional_beg",
	// additional tokens, handled in an ad-hoc manner
	"TILDE",
	"additional_end",
}

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
