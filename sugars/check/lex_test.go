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
	checkPos    int
	checkEnd    int
	identifiers []string
}

func lexParserNodeComparison(want lexParserNode, got Statement) (string, bool) {
	if want.completed != got.isCompleted {
		return fmt.Sprintf("got isCompleted=%v, want %v", got.isCompleted, want.completed), false
	}

	if msg, ok := lextest.CompareOptionalPos("pos", want.pos, got.pos); !ok {
		return msg, false
	}

	if msg, ok := lextest.CompareOptionalPos("end", want.end, got.end); !ok {
		return msg, false
	}

	if msg, ok := lextest.CompareOptionalPos("checkPos", want.checkPos, got.checkPos); !ok {
		return msg, false
	}

	if msg, ok := lextest.CompareOptionalPos("checkEnd", want.checkEnd, got.checkEnd); !ok {
		return msg, false
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
					completed: true,
					pos:       16, // check
					end:       41, // ;
					checkPos:  16, checkEnd: 22,
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
					completed: true,
					pos:       16, // check
					end:       41, // ;
					checkPos:  16, checkEnd: 22,
				},
				{
					completed: true,
					pos:       43, // check
					end:       62, // ;
					checkPos:  43, checkEnd: 49,
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
					completed: true,
					pos:       16, // x
					end:       40, // ;
					checkPos:  21, checkEnd: 27,
					identifiers: []string{"x"},
				},
			},
		},

		{
			name: "valid: 2 lines with different paths",
			code: makeCode(
				`func test() {`,
				`	check   path.Resolve("/")`,
				`	x := check svc.Field.DoSomething()`,
				`}`,
			),
			expected: []lexParserNode{
				{
					completed: true,
					pos:       16, // check
					end:       41, // ;
					checkPos:  16, checkEnd: 24,
				},

				{
					completed: true,
					pos:       43, // x
					end:       77, // ;
					checkPos:  48, checkEnd: 54,
					identifiers: []string{"x"},
				},
			},
		},

		{
			name: "whole codeblock",
			code: `package example

import (
	"fmt"
	"strconv"
)

func test() {
	check doSomething()
	x := check   strconv.Atoi("123")

	fmt.Println(x)
}
`,
			expected: []lexParserNode{
				{
					completed: true,
					pos:       63, // check
					end:       82, // ;
					checkPos:  63, checkEnd: 69,
				},

				{
					completed: true,
					pos:       84,  // x
					end:       116, // ;
					checkPos:  89, checkEnd: 97,
					identifiers: []string{"x"},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parser := LexicalParser()
			result := lextest.ExecuteLexicalParserContinuously(parser, tc.code, lextest.AsType[Statement])

			lextest.AssertNodes(t, tc.code, tc.expected, result, lexParserNodeComparison)
		})
	}
}
