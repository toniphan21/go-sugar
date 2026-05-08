package check

import (
	"testing"

	"nhatp.com/go/sugar/lextest"
)

var makeCode = lextest.MakeCode

type lexParserNode struct {
	completed   bool
	pos         int
	end         int
	operandPkg  string
	operandName string
}

func (w *lexParserNode) assert(t *testing.T, code string, got Statement, idx int) {
	t.Helper()
	valid := true

	want := w

	if want.completed != got.isCompleted {
		t.Errorf("index %d: got isCompleted=%v, want %v", idx, got.isCompleted, want.completed)
		valid = false
	}

	gotPos := -1
	if got.pos != nil {
		gotPos = int(got.pos.Pos)
	}
	if want.pos != gotPos {
		t.Errorf("index %d: got pos=%d, want %d", idx, gotPos, want.pos)
		valid = false
	}

	gotEnd := -1
	if got.end != nil {
		gotEnd = int(got.end.Pos)
	}
	if want.end != gotEnd {
		t.Errorf("index %d: got end=%d, want %d", idx, gotEnd, want.pos)
		valid = false
	}

	gOP := ""
	if got.operandPkg != nil {
		gOP = *got.operandPkg
	}
	if want.operandPkg != gOP {
		t.Errorf("index %d: got operandPkg=%v, want %v", idx, got.operandPkg, want.operandPkg)
		valid = false
	}

	if want.operandName != got.operandName {
		t.Errorf("index %d: got operandName=%v, want %v", idx, got.operandName, want.operandName)
		valid = false
	}

	if !valid {
		t.Log(lextest.LogMessageForLexViewer(code))
	}
}

type lexParserTestCase struct {
	name     string
	code     string
	expected []lexParserNode
}

func Test_Recognizer(t *testing.T) {
	cases := []lexParserTestCase{
		{
			name:     "empty",
			code:     "",
			expected: nil,
		},

		{
			name: "invalid: missing OperandName",
			code: makeCode(
				`func test() {`,
				`	check 123"`,
				`}`,
			),
			expected: nil,
		},

		{
			name: "valid: 1 line",
			code: makeCode(
				`func test() {`,
				`	check strconv.Atoi("123")`,
				`}`,
			),
			expected: []lexParserNode{
				{
					completed:   true,
					pos:         16, // check
					end:         34, // (
					operandPkg:  "strconv",
					operandName: "Atoi",
				},
			},
		},

		{
			name: "valid: 2 lines",
			code: makeCode(
				`func test() {`,
				`	check strconv.Atoi("123")`,
				`	check doSomething()`,
				`}`,
			),
			expected: []lexParserNode{
				{
					completed:   true,
					pos:         16, // check
					end:         34, // (
					operandPkg:  "strconv",
					operandName: "Atoi",
				},
				{
					completed:   true,
					pos:         43, // check
					end:         60, // (
					operandName: "doSomething",
				},
			},
		},

		{
			name: "valid: 1 line with IdentifierLHS",
			code: makeCode(
				`func test() {`,
				`	x := check doSomething()`,
				`}`,
			),
			expected: []lexParserNode{
				{
					completed:   true,
					pos:         16, // x
					end:         38, // (
					operandPkg:  "",
					operandName: "doSomething",
				},
			},
		},

		{
			name: "valid: 2 lines with different paths",
			code: makeCode(
				`func test() {`,
				`	check path.Resolve("/")`,
				`	x := check doSomething()`,
				`}`,
			),
			expected: []lexParserNode{
				{
					completed:   true,
					pos:         16, // check
					end:         34, // (
					operandPkg:  "path",
					operandName: "Resolve",
				},

				{
					completed:   true,
					pos:         41, // x
					end:         63, // (
					operandPkg:  "",
					operandName: "doSomething",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reg := LexicalParser()
			result := lextest.ExecuteLexicalParserContinuouslyWithCheckpoint(reg, tc.code, lextest.AsType[Statement])

			if len(result) != len(tc.expected) {
				t.Log(lextest.LogMessageForLexViewer(tc.code))
				t.Errorf("len(result) = %d, want %d", len(result), len(tc.expected))
			}

			for i, v := range result {
				tc.expected[i].assert(t, tc.code, v, i)
			}
		})
	}
}
