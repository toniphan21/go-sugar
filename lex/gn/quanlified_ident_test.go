package gn

import (
	"testing"

	"nhatp.com/go/sugar/lextest"
)

func qualifiedIdentTestCases() []lextest.ContinuousTestCase[qualifiedIdent] {
	return []lextest.ContinuousTestCase[TypeName]{
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
}

func runQualifiedIdentTest(t *testing.T, tc lextest.ContinuousTestCase[qualifiedIdent]) {
	t.Helper()

	parser := TypeNameParser()
	result := lextest.ExecuteLexicalParserContinuously(parser, tc.Code, lextest.AsType[TypeName])

	if len(result) != len(tc.Expected) {
		t.Errorf("len(result) = %d; want %d", len(result), len(tc.Expected))
	}

	for i, v := range tc.Expected {
		if tc.Expected[i].Identifier != result[i].Identifier {
			t.Errorf("result[%d].Identifier = %s; want %s", i, result[i].Identifier, v.Identifier)
		}

		wantPkg, gotPkg := "", ""
		if tc.Expected[i].PackageName != nil {
			wantPkg = *tc.Expected[i].PackageName
		}
		if result[i].PackageName != nil {
			gotPkg = *result[i].PackageName
		}
		if wantPkg != gotPkg {
			t.Errorf("result[i].PackageName = %s; want %s", gotPkg, wantPkg)
		}
	}
}

func Test_OperandName(t *testing.T) {
	cases := qualifiedIdentTestCases()
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runQualifiedIdentTest(t, tc)
		})
	}
}

func Test_TypeName(t *testing.T) {
	cases := qualifiedIdentTestCases()
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runQualifiedIdentTest(t, tc)
		})
	}
}
