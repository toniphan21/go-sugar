package lex

import (
	"fmt"
	"testing"

	"nhatp.com/go/sugar/lextest"
)

type callSuffixNode struct {
	pos int
	end int
}

func callSuffixNodeComparison(want callSuffixNode, got CallSuffix) (string, bool) {
	if want.pos != int(got.Pos.Pos) {
		return fmt.Sprintf("got pos=%d, want %d", int(got.Pos.Pos), want.pos), false
	}

	if want.end != int(got.End.Pos) {
		return fmt.Sprintf("got end=%d, want %d", int(got.End.Pos), want.end), false
	}
	return "", true
}

func callSuffixNodes(item ...callSuffixNode) []callSuffixNode {
	return item
}

func Test_CallSuffix(t *testing.T) {
	cases := []lextest.ContinuousTestCase[callSuffixNode]{
		{
			Name:     "invalid: single lit with nothing",
			Code:     `1`,
			Expected: nil,
		},

		{
			Name:     "invalid: func at start",
			Code:     `func`,
			Expected: nil,
		},

		{
			Name:     "invalid: ident start",
			Code:     `a`,
			Expected: nil,
		},

		{
			Name:     "invalid: never close",
			Code:     `(((aa))`,
			Expected: nil,
		},

		{
			Name:     "valid: simple",
			Code:     `()`,
			Expected: lextest.Nodes(callSuffixNode{pos: 1, end: 3}),
		},

		{
			Name:     "valid: one with ident and literal",
			Code:     `(a, true, "123")`,
			Expected: lextest.Nodes(callSuffixNode{pos: 1, end: 17}),
		},

		{
			Name:     "valid: first one failed, second one ok",
			Code:     "call(a\ncall(b)",
			Expected: lextest.Nodes(callSuffixNode{pos: 12, end: 15}),
		},

		{
			Name:     "valid: first one ok, second one fail",
			Code:     "call(a)\ncall(b\n",
			Expected: lextest.Nodes(callSuffixNode{pos: 5, end: 8}),
		},

		{
			Name:     "valid: 2 calls",
			Code:     "call(a)\ncall(b)\n",
			Expected: lextest.Nodes(callSuffixNode{pos: 5, end: 8}, callSuffixNode{pos: 13, end: 16}),
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			parser := CallSuffixParser()
			result := lextest.ExecuteLexicalParserContinuously(parser, tc.Code, lextest.AsType[CallSuffix])

			lextest.AssertNodes(t, tc.Code, tc.Expected, result, callSuffixNodeComparison)
		})
	}
}
