package require

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/lex"
)

func New() sugar.Sugar {
	return &sugarImpl{Base: &sugar.Base{}}
}

const keyword = "require"
const ID = "nhatp.com/go/sugar/sugars/require"
const NodeType = ID + ".node"
const Binary = "sugar-require"
const Version = sugar.Version

type sugarImpl struct {
	*sugar.Base
}

func (s *sugarImpl) ID() string {
	return ID
}

func (s *sugarImpl) Binary() sugar.BinaryInfo {
	return sugar.BinaryInfo{
		Name:    Binary,
		Version: Version,
		Usage:   "a go-sugar plugin for \"require\" sugar.\n\nSee full usage at https://nhatp.com/go/sugar/sugars/check",
	}
}

func (s *sugarImpl) Parse(source []byte) []sugar.Node {
	return sugar.DoParse(LexicalParser(), source)
}

func (s *sugarImpl) Restore(source []byte) []sugar.Node {
	return sugar.DoParse(lex.SugarPlaceholderFuncParser(keyword), source)
}

func (s *sugarImpl) StructuralTransform(sourceID string, n sugar.Node) ([]byte, error) {
	return sugar.DoTransform[*node](s.Base, sourceID, n, func(source []byte, data *node) ([]byte, error) {
		out := bytes.Buffer{}

		out.Write(source[data.pos.Offset:data.requirePos.Offset])
		out.Write([]byte(lex.SugarPlaceholderFuncName(keyword)))
		out.WriteRune('(')
		if data.message == nil {
			out.Write(source[data.requireEnd.Offset:data.end.Offset])
		} else {
			out.Write(source[data.requireEnd.Offset:data.messagePos.Offset])
			out.WriteRune(',')
			out.Write(source[data.messagePos.Offset:data.end.Offset])
		}
		out.WriteRune(')')

		return out.Bytes(), nil
	})
}

func (s *sugarImpl) SemanticTransformer(sourceID string, scopeID string, n sugar.Node) ([]byte, error) {
	return sugar.DoTransformWithSemanticScope[*node](s.Base, sourceID, scopeID, n, func(source []byte, data *node) ([]byte, error) {
		var expr []byte
		if data.messagePos == nil {
			expr = source[data.requireEnd.Offset:data.end.Offset]
		} else {
			expr = source[data.requireEnd.Offset:data.messagePos.Offset]
		}

		out := bytes.Buffer{}

		if len(data.identifiers) == 0 {
			out.WriteString("if err := ")
			out.Write(expr)
			out.WriteString("; err != nil {\n")
		} else {
			idents := make([]string, len(data.identifiers)+1)
			copy(idents, data.identifiers)
			idents[len(data.identifiers)] = "err"

			out.WriteString(strings.Join(idents, ", "))
			out.WriteString(" := ")
			out.Write(expr)
			out.WriteRune('\n')
			out.WriteString("if err != nil {\n")
		}

		out.WriteRune('\t')
		s.emitHandleCode(&out, data, expr)
		out.WriteRune('}')

		return out.Bytes(), nil
	})
}

func (s *sugarImpl) emitHandleCode(out *bytes.Buffer, data *node, expr []byte) {
	name, have := data.findTestingParam()
	if !have {
		out.WriteString(`panic(`)
		if data.message == nil {
			out.WriteString(strconv.Quote(strings.TrimSpace(string(expr)) + ": "))
			out.WriteString(" + err.Error())\n")
			return
		}

		out.WriteString(*data.message)
		out.WriteString(")\n")
		return
	}

	out.WriteString(name)

	exprString := strconv.Quote(strings.TrimSpace(string(expr)))
	if data.message == nil {
		out.WriteString(`.Fatalf("%s: %w", `)
		out.WriteString(exprString)
		out.WriteString(", err)\n")
		return
	}

	msg := *data.message
	switch hasS, hasW := data.scanMessageVerbs(msg); {
	case hasS && hasW:
		out.WriteString(`.Fatalf(`)
		out.WriteString(msg)
		out.WriteString(`, `)
		out.WriteString(exprString)
		out.WriteString(", err)\n")
	case hasS:
		out.WriteString(`.Fatalf(`)
		out.WriteString(msg)
		out.WriteString(`, `)
		out.WriteString(exprString)
		out.WriteString(")\n")
	case hasW:
		out.WriteString(`.Fatalf(`)
		out.WriteString(msg)
		out.WriteString(", err)\n")
	default:
		out.WriteString(`.Fatal(`)
		out.WriteString(msg)
		out.WriteString(")\n")
	}
}

func (s *sugarImpl) RestoreTransform(sourceID string, n sugar.Node) ([]byte, error) {
	return sugar.DoTransform[lex.SugarPlaceholderFunc](s.Base, sourceID, n, func(source []byte, data lex.SugarPlaceholderFunc) ([]byte, error) {
		if data.Keyword() != keyword {
			return nil, sugar.ErrUnknownNode
		}

		body := data.Body()
		lb := len(body)
		if lb == 0 {
			return nil, fmt.Errorf("empty body")
		}

		is := &sugar.LexemePredicate{}
		if lb > 2 && is.String(body[lb-1]) && is.Comma(body[lb-2]) {
			// end with message, we need to strip the comma before that
			out := bytes.Buffer{}
			out.WriteString("require ")
			out.Write(source[data.InnerPos().Offset:body[lb-2].Offset]) // inner pos til the comma
			out.WriteRune(' ')
			out.Write(source[body[lb-1].Offset:data.InnerEnd().Offset]) // literal string after comma

			return out.Bytes(), nil
		}

		// no message, just unwrap the function
		out := bytes.Buffer{}
		out.WriteString("require ")
		out.Write(source[data.InnerPos().Offset:data.InnerEnd().Offset])

		return out.Bytes(), nil
	})
}
