package ebnf

import (
	"reflect"
	"testing"

	"nhatp.com/go/sugar/lextest"
)

func Test_OperandName(t *testing.T) {
	cases := []lextest.ContinuousTestCase[OperandName]{
		{
			Name:     "not found",
			Code:     `"hello"`,
			Expected: nil,
		},

		{
			Name:     "invalid: missing identifier",
			Code:     `pkg.`,
			Expected: nil,
		},

		{
			Name: "valid: x.y",
			Code: `x.y`,
			Expected: []OperandName{
				{PackageName: new("x"), Identifier: "y"},
			},
		},

		{
			Name: "valid: 2 positions",
			Code: `x := strconv.Atoi("1")`,
			Expected: []OperandName{
				{Identifier: "x"},
				{PackageName: new("strconv"), Identifier: "Atoi"},
			},
		},

		{
			Name: "valid: at . it's stop then restart at Name",
			Code: `.Name`,
			Expected: []OperandName{
				{Identifier: "Name"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			parser := OperandNameParser()
			result := lextest.ExecuteLexicalParserContinuously(parser, tc.Code, lextest.AsType[OperandName])

			if !reflect.DeepEqual(result, tc.Expected) {
				t.Errorf("expected %v but got %v", tc.Expected, result)
			}
		})
	}
}
