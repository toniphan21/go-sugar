package ebnf

import (
	"reflect"
	"testing"

	"nhatp.com/go/sugar/lextest"
)

func Test_IdentifierList(t *testing.T) {
	cases := []lextest.ContinuousTestCase[IdentifierList]{
		{
			Name:     "invalid: single lit",
			Code:     `"hello"`,
			Expected: nil,
		},

		{
			Name:     "invalid: multiple literals",
			Code:     `1, 2, "a", 'b'`,
			Expected: nil,
		},

		{
			Name:     "invalid: mixed, stop before lit",
			Code:     `a, 1`,
			Expected: nil,
		},

		{
			Name: "valid: single ident",
			Code: `a`,
			Expected: []IdentifierList{
				{Identifier: []string{"a"}},
			},
		},

		{
			Name: "valid: two idents",
			Code: `a, b`,
			Expected: []IdentifierList{
				{Identifier: []string{"a", "b"}},
			},
		},

		{
			Name: "valid: restart every time see .",
			Code: `t.a, t.b`,
			Expected: []IdentifierList{
				{Identifier: []string{"t"}},
				{Identifier: []string{"a", "t"}},
				{Identifier: []string{"b"}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			parser := IdentifierListParser()
			result := lextest.ExecuteLexicalParserContinuously(parser, tc.Code, lextest.AsType[IdentifierList])

			if !reflect.DeepEqual(result, tc.Expected) {
				t.Errorf("expected %v but got %v", tc.Expected, result)
			}
		})
	}
}
