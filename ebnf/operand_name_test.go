package ebnf

import (
	"reflect"
	"testing"
)

func Test_OperandName(t *testing.T) {
	cases := []lexicalParserTestCase[OperandName]{
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
			expected: []OperandName{
				{PackageName: new("x"), Identifier: "y"},
			},
		},

		{
			name: "valid: 2 positions",
			code: `x := strconv.Atoi("1")`,
			expected: []OperandName{
				{Identifier: "x"},
				{PackageName: new("strconv"), Identifier: "Atoi"},
			},
		},

		{
			name: "valid: at . it's stop then restart at Name",
			code: `.Name`,
			expected: []OperandName{
				{Identifier: "Name"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parser := OperandNameParser()
			result := executeLexicalParserContinuously(parser, tc.code, asType[OperandName])

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("expected %v but got %v", tc.expected, result)
			}
		})
	}
}
