package main

import (
	"encoding/json"
	"fmt"

	"nhatp.com/go/sugar"
)

func safeExtractLexemes(input string) (output string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	return extractLexemes(input)
}

type Lexeme struct {
	Token           string          `json:"Tok"`
	LexemePredicate LexemePredicate `json:"LexemePredicate"`

	Lit    string `json:"lit"`
	Raw    string `json:"raw"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
	Offset int    `json:"offset"`
}

type LexemePredicate struct {
	StatementBoundary bool `json:"StatementBoundary"`
}

func extractLexemes(input string) (string, error) {
	content := []byte(input)
	lex := sugar.Lex(content)
	var list []*Lexeme

	on := &sugar.LexemePredicate{}

	for i, v := range lex {
		end := len(content)
		if i+1 < len(lex) {
			end = lex[i+1].Offset
		}

		item := &Lexeme{
			Lit:    v.Lit,
			Line:   v.Line,
			Column: v.Column,
			Offset: v.Offset,
			Raw:    string(content[v.Offset:end]),
			LexemePredicate: LexemePredicate{
				StatementBoundary: on.Boundary(v),
			},
		}

		t := int(v.Tok)
		if t < len(tokenNames) {
			item.Token = tokenNames[t]
		} else {
			item.Token = v.Tok.String()
		}

		list = append(list, item)
	}

	out, err := json.Marshal(list)
	if err != nil {
		return "[]", err
	}
	return string(out), nil
}
