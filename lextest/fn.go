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
	for _, v := range lexemes {
		if parser.Done(v) {
			item, ok := parser.Result()
			if ok {
				node, ok := as(item)
				if ok {
					result = append(result, node)
				}
			}
			parser.Reset()
		}
	}
	return result
}

func FormatCodeForLexViewer(code string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(code))

	return fmt.Sprintf("code:base64:%s", encoded)
}

func LogMessageForLexViewer(code string) string {
	return fmt.Sprintf("copy line below to tools/lexeme-viewer for debugging\n%v\n", FormatCodeForLexViewer(code))
}
