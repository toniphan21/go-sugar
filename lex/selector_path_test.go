package lex

import (
	"reflect"
	"testing"

	"nhatp.com/go/sugar/lextest"
)

func Test_SelectorPath(t *testing.T) {
	cases := []lextest.ContinuousTestCase[SelectorPath]{
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
			Name:     "invalid: ident dot literal",
			Code:     `a."1"`,
			Expected: nil,
		},

		{
			Name: "valid: single ident",
			Code: `a`,
			Expected: []SelectorPath{
				{Identifiers: []string{"a"}},
			},
		},

		{
			Name: "valid: two idents",
			Code: `a.b`,
			Expected: []SelectorPath{
				{Identifiers: []string{"a", "b"}},
			},
		},

		{
			Name: "valid: two paths",
			Code: "a.b\nc.d.e = 1",
			Expected: []SelectorPath{
				{Identifiers: []string{"a", "b"}},
				{Identifiers: []string{"c", "d", "e"}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			parser := SelectorPathParser()
			result := lextest.ExecuteLexicalParserContinuously(parser, tc.Code, lextest.AsType[SelectorPath])

			if !reflect.DeepEqual(result, tc.Expected) {
				t.Errorf("expected %v but got %v", tc.Expected, result)
			}
		})
	}
}
