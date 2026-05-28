package lex

import (
	"fmt"
	"testing"

	"nhatp.com/go/sugar/lextest"
)

type sugarPlaceholderFuncNode struct {
	pos      int
	end      int
	innerPos int
	innerEnd int
}

func sugarPlaceholderFuncNodeComparison(want sugarPlaceholderFuncNode, got SugarPlaceholderFunc) (string, bool) {
	if want.pos != int(got.pos.Pos) {
		return fmt.Sprintf("got pos=%d, want %d", int(got.pos.Pos), want.pos), false
	}

	if want.end != int(got.end.Pos) {
		return fmt.Sprintf("got pos=%d, want %d", int(got.end.Pos), want.end), false
	}

	if want.innerPos != int(got.innerPos.Pos) {
		return fmt.Sprintf("got pos=%d, want %d", int(got.innerPos.Pos), want.innerPos), false
	}

	if want.innerEnd != int(got.innerEnd.Pos) {
		return fmt.Sprintf("got pos=%d, want %d", int(got.innerEnd.Pos), want.innerEnd), false
	}

	return "", true
}

func Test_SugarFunc(t *testing.T) {
	cases := []lextest.ContinuousTestCase[sugarPlaceholderFuncNode]{
		{
			Name:     "invalid: single lit with nothing",
			Code:     `1`,
			Expected: nil,
		},

		{
			Name:     "invalid: func with nothing",
			Code:     `__sugar_test__`,
			Expected: nil,
		},

		{
			Name:     "invalid: func with empty call",
			Code:     `__sugar_test__()`,
			Expected: nil,
		},

		{
			Name:     "invalid: func with open without close",
			Code:     `__sugar_test__(123`,
			Expected: nil,
		},

		{
			Name:     "valid: wrap literal",
			Code:     `__sugar_test__(123)`,
			Expected: lextest.Nodes(sugarPlaceholderFuncNode{pos: 1, end: 20, innerPos: 16, innerEnd: 19}),
		},

		{
			Name:     "valid: wrap function call",
			Code:     `__sugar_test__(doSomething(""))`,
			Expected: lextest.Nodes(sugarPlaceholderFuncNode{pos: 1, end: 32, innerPos: 16, innerEnd: 31}),
		},

		{
			Name: "valid: wrap function call with multi-line",
			Code: `func test() {
x := 1
__sugar_test__(doSomething("any", x))
}`,
			Expected: lextest.Nodes(sugarPlaceholderFuncNode{pos: 22, end: 59, innerPos: 37, innerEnd: 58}),
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			parser := SugarPlaceholderFuncParser("test")
			result := lextest.ExecuteLexicalParserContinuously(parser, tc.Code, lextest.AsType[SugarPlaceholderFunc])

			lextest.AssertNodes(t, tc.Code, tc.Expected, result, sugarPlaceholderFuncNodeComparison)
		})
	}
}
