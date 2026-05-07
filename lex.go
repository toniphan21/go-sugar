package sugar

import (
	"fmt"
	"go/scanner"
	"go/token"
)

type Lexeme struct {
	Tok    token.Token
	Pos    token.Pos
	Lit    string
	Line   int
	Column int
	Offset int
}

func (l *Lexeme) GoString() string {
	return fmt.Sprintf("%v | Tok:%d Pos:%d Lit:%#v | line:%d col:%d offset:%d", TokenName(l.Tok), l.Tok, l.Pos, l.Lit, l.Line, l.Column, l.Offset)
}

func Lex(content []byte) []Lexeme {
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(content))

	var s scanner.Scanner
	s.Init(file, content, nil, 0)

	var result []Lexeme
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}

		position := fset.Position(pos)
		lex := Lexeme{
			Tok:    tok,
			Pos:    pos,
			Lit:    lit,
			Line:   position.Line,
			Column: position.Column,
			Offset: position.Offset,
		}
		if lex.Lit == "" {
			lex.Lit = tok.String()
		}
		result = append(result, lex)
	}
	return result
}

// ---

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
