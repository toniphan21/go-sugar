package check

import (
	"fmt"
	"reflect"
	"testing"

	"nhatp.com/go/sugar/lextest"
)

var makeCode = lextest.MakeCode

type lexParserNode struct {
	completed   bool
	pos         int
	end         int
	identifiers []string
}

func lexParserNodeComparison(want lexParserNode, got Statement) (string, bool) {
	if want.completed != got.isCompleted {
		return fmt.Sprintf("got isCompleted=%v, want %v", got.isCompleted, want.completed), false
	}

	gotPos := -1
	if got.pos != nil {
		gotPos = int(got.pos.Pos)
	}
	if want.pos != gotPos {
		return fmt.Sprintf("got pos=%d, want %d", gotPos, want.pos), false
	}

	gotEnd := -1
	if got.end != nil {
		gotEnd = int(got.end.Pos)
	}
	if want.end != gotEnd {
		return fmt.Sprintf("got end=%d, want %d", gotEnd, want.pos), false
	}

	if !reflect.DeepEqual(want.identifiers, got.identifiers) {
		return fmt.Sprintf("got identifiers=%v, want %v", got.identifiers, want.identifiers), false
	}
	return "", true
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
					end:         41, // ;
					identifiers: []string{"strconv", "Atoi"},
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
					end:         41, // ;
					identifiers: []string{"strconv", "Atoi"},
				},
				{
					completed:   true,
					pos:         43, // check
					end:         62, // ;
					identifiers: []string{"doSomething"},
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
					end:         40, // ;
					identifiers: []string{"doSomething"},
				},
			},
		},

		{
			name: "valid: 2 lines with different paths",
			code: makeCode(
				`func test() {`,
				`	check path.Resolve("/")`,
				`	x := check svc.Field.DoSomething()`,
				`}`,
			),
			expected: []lexParserNode{
				{
					completed:   true,
					pos:         16, // check
					end:         39, // ;
					identifiers: []string{"path", "Resolve"},
				},

				{
					completed:   true,
					pos:         41, // x
					end:         75, // ;
					identifiers: []string{"svc", "Field", "DoSomething"},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reg := LexicalParser()
			result := lextest.ExecuteLexicalParserContinuously(reg, tc.code, lextest.AsType[Statement])

			lextest.AssertNodes(t, tc.code, tc.expected, result, lexParserNodeComparison)
		})
	}
}
