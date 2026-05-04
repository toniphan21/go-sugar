package parse

import (
	"reflect"
	"strings"
	"testing"

	"nhatp.com/go/sugar"
)

func makeCode(lines ...string) string {
	return strings.Join(lines, "\n")
}

func Test_OperandName(t *testing.T) {
	cases := []struct {
		name     string
		code     string
		expected []*OperandName
	}{
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
			expected: []*OperandName{
				{PackageName: new("x"), Identifier: "y"},
			},
		},

		{
			name: "valid: 2 positions",
			code: `x := strconv.Atoi("1")`,
			expected: []*OperandName{
				{Identifier: "x"},
				{PackageName: new("strconv"), Identifier: "Atoi"},
			},
		},

		{
			name: "valid: at . it's stop then restart at Name",
			code: `.Name`,
			expected: []*OperandName{
				{Identifier: "Name"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			lexemes := sugar.Lex([]byte(tc.code))

			var result []*OperandName
			parser := OperandNameLexicalParser()
			for _, v := range lexemes {
				switch parser.Take(v) {
				case sugar.StatusCompleted:
					result = append(result, parser.Build())
					parser.Reset()

				case sugar.StatusFailed:
					parser.Reset()

				default:
					// keep go-ing
				}
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("expected %v but got %v", tc.expected, result)
			}
		})
	}
}
