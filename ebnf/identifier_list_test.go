package ebnf

import (
	"reflect"
	"testing"
)

func Test_IdentifierList(t *testing.T) {
	cases := []lexicalParserTestCase[IdentifierList]{
		{
			name:     "invalid: single lit",
			code:     `"hello"`,
			expected: nil,
		},

		{
			name:     "invalid: multiple literals",
			code:     `1, 2, "a", 'b'`,
			expected: nil,
		},

		{
			name:     "invalid: mixed, stop before lit",
			code:     `a, 1`,
			expected: nil,
		},

		{
			name: "valid: single ident",
			code: `a`,
			expected: []IdentifierList{
				{Identifier: []string{"a"}},
			},
		},

		{
			name: "valid: two idents",
			code: `a, b`,
			expected: []IdentifierList{
				{Identifier: []string{"a", "b"}},
			},
		},

		{
			name: "valid: restart every time see .",
			code: `t.a, t.b`,
			expected: []IdentifierList{
				{Identifier: []string{"t"}},
				{Identifier: []string{"a", "t"}},
				{Identifier: []string{"b"}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parser := IdentifierListParser()
			result := executeLexicalParserContinuously(parser, tc.code, asType[IdentifierList])

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("expected %v but got %v", tc.expected, result)
			}
		})
	}
}
