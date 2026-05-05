package ebnf

import (
	"reflect"
	"testing"
)

func Test_TypeName(t *testing.T) {
	cases := []lexicalParserTestCase[TypeName]{
		{
			name:     "not found",
			code:     `"hello"`,
			expected: nil,
		},

		{
			name:     "invalid: missing identifier",
			code:     `pkg.`,
			expected: nil,
		},

		{
			name: "valid: x.y",
			code: `x.y`,
			expected: []TypeName{
				{PackageName: new("x"), Identifier: "y"},
			},
		},

		{
			name: "valid: 2 positions",
			code: `x := strconv.Atoi("1")`,
			expected: []TypeName{
				{Identifier: "x"},
				{PackageName: new("strconv"), Identifier: "Atoi"},
			},
		},

		{
			name: "valid: at . it's stop then restart at Name",
			code: `.Name`,
			expected: []TypeName{
				{Identifier: "Name"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parser := TypeNameParser()
			result := executeLexicalParserContinuously(parser, tc.code, asType[TypeName])

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("expected %v but got %v", tc.expected, result)
			}
		})
	}
}
