package ebnf

import (
	"reflect"
	"testing"

	"nhatp.com/go/sugar/lextest"
)

func Test_TypeName(t *testing.T) {
	cases := []lextest.ContinuousTestCase[TypeName]{
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
			Expected: []TypeName{
				{PackageName: new("x"), Identifier: "y"},
			},
		},

		{
			Name: "valid: 2 positions",
			Code: `x := strconv.Atoi("1")`,
			Expected: []TypeName{
				{Identifier: "x"},
				{PackageName: new("strconv"), Identifier: "Atoi"},
			},
		},

		{
			Name: "valid: at . it's stop then restart at Name",
			Code: `.Name`,
			Expected: []TypeName{
				{Identifier: "Name"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			parser := TypeNameParser()
			result := lextest.ExecuteLexicalParserContinuously(parser, tc.Code, lextest.AsType[TypeName])

			if !reflect.DeepEqual(result, tc.Expected) {
				t.Errorf("expected %v but got %v", tc.Expected, result)
			}
		})
	}
}
