package lex

import (
	"reflect"
	"testing"

	"nhatp.com/go/sugar/lextest"
)

func Test_IdentifierLHS(t *testing.T) {
	cases := []lextest.ContinuousTestCase[IdentifierLHS]{
		{
			Name:     "invalid: single ident with nothing",
			Code:     `a`,
			Expected: nil,
		},

		{
			Name:     "invalid: single ident with assign",
			Code:     `a =`,
			Expected: nil,
		},

		{
			Name:     "invalid: two ident end with literal",
			Code:     `a, b, "1"`,
			Expected: nil,
		},

		{
			Name: "valid: single ident with define",
			Code: `a :=`,
			Expected: []IdentifierLHS{
				{Identifiers: []string{"a"}},
			},
		},

		{
			Name: "valid: discard first line starts with literal",
			Code: "1\na :=",
			Expected: []IdentifierLHS{
				{Identifiers: []string{"a"}},
			},
		},

		{
			Name: "valid: 2 idents with define",
			Code: `a, b :=`,
			Expected: []IdentifierLHS{
				{Identifiers: []string{"a", "b"}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			parser := IdentifierLHSParser()
			result := lextest.ExecuteLexicalParserContinuously(parser, tc.Code, lextest.AsType[IdentifierLHS])

			if !reflect.DeepEqual(result, tc.Expected) {
				t.Errorf("expected %v but got %v", tc.Expected, result)
			}
		})
	}
}
