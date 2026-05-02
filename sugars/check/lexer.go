package check

import (
	"go/token"
	"strings"

	"nhatp.com/go/sugar"
)

type checkToken struct {
	original *sugar.Token
	lhs      []token.Token
	operand  string
}

func (c *checkToken) Tok() token.Token {
	return c.original.Tok()
}

func (c *checkToken) Pos() token.Pos {
	return c.original.Pos()
}

func (c *checkToken) Lit() string {
	return c.original.Lit()
}

func (c *checkToken) Offset() int {
	return c.original.Offset()
}

func (c *checkToken) Raw() string {
	return strings.Repeat(" ", len(c.original.Raw()))
}

var _ sugar.Lexeme = (*checkToken)(nil)

func Scan(source []byte) []sugar.Lexeme {
	return sugar.Scan(source, func(prev sugar.Lexeme, current *sugar.Token, next *sugar.Token, source []byte) sugar.Lexeme {
		if current.Lit() == "check" {
			return &checkToken{
				original: current,
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
