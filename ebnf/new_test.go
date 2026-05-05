package ebnf

import (
	"strings"

	"nhatp.com/go/sugar"
)

func asType[T any](in any) (T, bool) {
	v, ok := in.(T)
	return v, ok
}

func makeCode(lines ...string) string {
	return strings.Join(lines, "\n")
}

type lexicalParserTestCase[N any] struct {
	name     string
	code     string
	expected []N
}

func executeLexicalParserContinuously[P sugar.LexicalParser, N any](parser P, code string, as func(in any) (N, bool)) []N {
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
