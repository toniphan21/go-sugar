package require

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
	requirePos  int
	requireEnd  int
	messagePos  int
	messageEnd  int
	message     *string
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

	if msg, ok := lextest.CompareOptionalPos("requirePos", want.requirePos, got.requirePos); !ok {
		return msg, false
	}

	if msg, ok := lextest.CompareOptionalPos("requireEnd", want.requireEnd, got.requireEnd); !ok {
		return msg, false
	}

	if msg, ok := lextest.CompareOptionalPos("messagePos", want.messagePos, got.messagePos); !ok {
		return msg, false
	}

	if msg, ok := lextest.CompareOptionalPos("messageEnd", want.messageEnd, got.messageEnd); !ok {
		return msg, false
	}

	if msg, ok := lextest.CompareOptionalStrings("message", want.message, got.message); !ok {
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
				`	require 123"`,
				`}`,
			),
			expected: nil,
		},

		{
			name: "valid: 1 line",
			code: makeCode(
				`func test() {`,
				`	require strconv.Atoi("123")`,
				`}`,
			),
			expected: []lexParserNode{
				{
					completed:  true,
					pos:        16, // require
					end:        43, // ;
					requirePos: 16, requireEnd: 24,
				},
			},
		},

		{
			name: "valid: 2 lines",
			code: makeCode(
				`func test() {`,
				`	require strconv.Atoi("123")`,
				`	require doSomething()`,
				`}`,
			),
			expected: []lexParserNode{
				{
					completed:  true,
					pos:        16, // require
					end:        43, // ;
					requirePos: 16, requireEnd: 24,
				},
				{
					completed:  true,
					pos:        45, // require
					end:        66, // ;
					requirePos: 45, requireEnd: 53,
				},
			},
		},

		{
			name: "valid: 1 line with IdentifierLHS",
			code: makeCode(
				`func test() {`,
				`	x := require doSomething()`,
				`}`,
			),
			expected: []lexParserNode{
				{
					completed:  true,
					pos:        16, // x
					end:        42, // ;
					requirePos: 21, requireEnd: 29,
					identifiers: []string{"x"},
				},
			},
		},

		{
			name: "valid: 1 line with optional message",
			code: makeCode(
				`func test() {`,
				`	require doSomething() "cannot doSomething"`,
				`}`,
			),
			expected: []lexParserNode{
				{
					completed:  true,
					pos:        16, // require
					end:        58, // ;
					requirePos: 16, requireEnd: 24,
					messagePos: 38, messageEnd: 58,
					message: new(`"cannot doSomething"`),
				},
			},
		},

		{
			name: "valid: 2 lines with different paths",
			code: makeCode(
				`func test() {`,
				`	require   path.Resolve("/")`,
				`	x := require svc.Field.DoSomething()`,
				`}`,
			),
			expected: []lexParserNode{
				{
					completed:  true,
					pos:        16, // require
					end:        43, // ;
					requirePos: 16, requireEnd: 26,
				},

				{
					completed:  true,
					pos:        45, // x
					end:        81, // ;
					requirePos: 50, requireEnd: 58,
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
	require doSomething()
	x := require   strconv.Atoi("123")

	fmt.Println(x)
}
`,
			expected: []lexParserNode{
				{
					completed:  true,
					pos:        63, // require
					end:        84, // ;
					requirePos: 63, requireEnd: 71,
				},

				{
					completed:  true,
					pos:        86,  // x
					end:        120, // ;
					requirePos: 91, requireEnd: 101,
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
