package main

import (
	"encoding/json"
	"fmt"
	"go/token"
	"reflect"

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
	Pos             token.Pos       `json:"Pos"`
	LexemePredicate map[string]bool `json:"LexemePredicate"`

	Lit    string `json:"lit"`
	Raw    string `json:"raw"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
	Offset int    `json:"offset"`
}

func makeLexemePredicate(v sugar.Lexeme) map[string]bool {
	on := &sugar.LexemePredicate{}
	onType := reflect.TypeOf(on)
	onVal := reflect.ValueOf(on)
	boolType := reflect.TypeOf(true)
	lexemeType := reflect.TypeOf(v)

	results := make(map[string]bool)

	for i := range onType.NumMethod() {
		m := onType.Method(i)
		mt := m.Type

		if mt.NumIn() != 2 || mt.NumOut() != 1 {
			continue
		}
		if mt.In(1) != lexemeType || mt.Out(0) != boolType {
			continue
		}

		out := onVal.Method(i).Call([]reflect.Value{reflect.ValueOf(v)})
		results[m.Name] = out[0].Bool()
	}
	return results
}

func extractLexemes(input string) (string, error) {
	content := []byte(input)
	lex := sugar.Lex(content)
	var list []*Lexeme

	for i, v := range lex {
		end := len(content)
		if i+1 < len(lex) {
			end = lex[i+1].Offset
		}

		item := &Lexeme{
			Pos:             v.Pos,
			Lit:             v.Lit,
			Line:            v.Line,
			Column:          v.Column,
			Offset:          v.Offset,
			Raw:             string(content[v.Offset:end]),
			LexemePredicate: makeLexemePredicate(v),
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
