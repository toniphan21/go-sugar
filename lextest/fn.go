package lextest

import (
	"encoding/base64"
	"fmt"
	"strings"

	"nhatp.com/go/sugar"
)

func MakeCode(lines ...string) string {
	return strings.Join(lines, "\n")
}

func AsType[T any](in any) (T, bool) {
	v, ok := in.(T)
	return v, ok
}

type ContinuousTestCase[N any] struct {
	Name     string
	Code     string
	Expected []N
}

func ExecuteLexicalParserContinuously[P sugar.LexicalParser, N any](parser P, code string, as func(in any) (N, bool)) []N {
	lexemes := sugar.Lex([]byte(code))

	var result []N
	offset := 0
	for offset < len(lexemes) {
		slice := lexemes[offset:]

		if parser.Done(slice) {
			item, ok := parser.Result()
			if ok {
				node, valid := as(item)
				if valid {
					result = append(result, node)
				}
			}
		}
		offset += parser.Consumed()
		parser.Reset()
	}
	return result
}

func FormatCodeForLexViewer(code string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(code))

	return fmt.Sprintf("code:base64:%s", encoded)
}
