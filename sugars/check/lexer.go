package check

import (
	"go/token"
	"strings"

	"nhatp.com/go/sugar"
)

type checkToken struct {
	lex     *sugar.Token
	lhs     []token.Token
	operand string
}

func (c *checkToken) Tok() token.Token {
	return c.lex.Tok()
}

func (c *checkToken) Lit() string {
	return c.lex.Lit()
}

func (c *checkToken) Raw() string {
	return c.lex.Raw()
}

func (c *checkToken) Line() int {
	return c.lex.Line()
}

func (c *checkToken) Column() int {
	return c.lex.Column()
}

func (c *checkToken) Offset() int {
	return c.lex.Offset()
}

var _ sugar.Lexeme = (*checkToken)(nil)

func Scan(source []byte) []sugar.Lexeme {
	return sugar.Scan(source, func(prev sugar.Lexeme, current *sugar.Token, next *sugar.Token, source []byte) sugar.Lexeme {
		if current.Lit() == "check" {
			return &checkToken{
				lex: current,
			}
		}
		return current
	})
}

func Join(tokens []sugar.Lexeme) string {
	var sb strings.Builder
	for _, t := range tokens {
		sb.WriteString(t.Raw())
	}
	return sb.String()
}
