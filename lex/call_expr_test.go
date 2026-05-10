package lex

import (
	"fmt"
	"reflect"
	"testing"

	"nhatp.com/go/sugar/lextest"
)

type callExprNode struct {
	identifiers []string
	callPos     int
	callEnd     int
}

func callExprNodeComparison(want callExprNode, got CallExpr) (string, bool) {
	if !reflect.DeepEqual(want.identifiers, got.Identifiers) {
		return fmt.Sprintf("got identifier=%#v, want %#v", got.Identifiers, want.identifiers), false
	}

	if want.callPos != int(got.CallPos.Pos) {
		return fmt.Sprintf("got pos=%d, want %d", int(got.CallPos.Pos), want.callPos), false
	}

	if want.callEnd != int(got.CallEnd.Pos) {
		return fmt.Sprintf("got end=%d, want %d", int(got.CallEnd.Pos), want.callEnd), false
	}
	return "", true
}

func callExprNodes(item ...callExprNode) []callExprNode {
	return item
}

func Test_CallExpr(t *testing.T) {
	cases := []lextest.ContinuousTestCase[callExprNode]{
		{Name: "invalid: single lit with nothing", Code: `1`},
		{Name: "invalid: ident without call", Code: `a.b.c`},
		{Name: "invalid: ident with half-call", Code: `a.b.c(()`},
		{
			Name: "valid: single ident call with lit",
			Code: `do("any")`,
			Expected: callExprNodes(callExprNode{
				identifiers: []string{"do"},
				callPos:     3, callEnd: 10,
			}),
		},
		{
			Name: "valid: multiple idents call with lit",
			Code: lextest.MakeCode(
				`var test = func() {`,
				`	x := 0`,
				`	strconv.Atoi("any")`,
				`}`,
			),
			Expected: callExprNodes(callExprNode{
				identifiers: []string{"strconv", "Atoi"},
				callPos:     42, callEnd: 49,
			}),
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			parser := CallExprParser()
			result := lextest.ExecuteLexicalParserContinuously(parser, tc.Code, lextest.AsType[CallExpr])

			lextest.AssertNodes(t, tc.Code, tc.Expected, result, callExprNodeComparison)
		})
	}
}
