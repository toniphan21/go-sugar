package lex

import (
	"encoding/json"
	"fmt"

	"github.com/oklog/ulid/v2"
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/internal/sdk"
	"nhatp.com/go/sugar/internal/sdk/transport"
)

/*
Diagram(
  Start({type:'complex'}),
  Stack('__sugar_[keyword]__'),
  Stack('('),
  Stack('*'),
  ZeroOrMore('*'),
  Stack(')'),
  Stack(';'),
  End({type:'complex'})
)
*/

/*
Diagram(
  Start({type:'complex'}),
  Stack(NonTerminal('doCollectPos'), '__sugar_[keyword]__'),
  Stack('('),
  Stack(NonTerminal('doCollectInnerPos'), '*'),
  ZeroOrMore('*'),
  Stack(NonTerminal('doCollectInnerEnd'), ')'),
  Stack(NonTerminal('doCollectEnd'), ';'),
  End({type:'complex'})
)
*/

func init() {
	transport.RegisterNodeDeserializer(SugarPlaceholderFuncID, deserializeSugarPlaceholderFunc)
}

func serializeSugarPlaceholderFunc(n *SugarPlaceholderFunc) (*sdk.Node, error) {
	var body []sdk.Lex
	for _, b := range n.body {
		body = append(body, b.ToSDKLex())
	}

	payload, err := json.Marshal(&sugarPlaceholderFuncPayload{
		InnerPos: n.innerPos.ToSDKLex(),
		InnerEnd: n.innerEnd.ToSDKLex(),
		Body:     body,
		Keyword:  n.keyword,
	})
	if err != nil {
		return nil, err
	}

	return &sdk.Node{
		ID:      ulid.Make().String(),
		Type:    SugarPlaceholderFuncID,
		Pos:     n.pos.ToSDKLex(),
		End:     n.end.ToSDKLex(),
		Payload: payload,
	}, nil
}

func deserializeSugarPlaceholderFunc(in sdk.Node) (sugar.Node, error) {
	n := &SugarPlaceholderFunc{}
	n.pos = sugar.FromSDKLex(in.Pos)
	n.end = sugar.FromSDKLex(in.End)

	var payload sugarPlaceholderFuncPayload
	if err := json.Unmarshal(in.Payload, &payload); err != nil {
		return nil, err
	}

	var body []sugar.Lexeme
	for _, b := range payload.Body {
		body = append(body, sugar.FromSDKLex(b))
	}

	n.innerPos = sugar.FromSDKLex(payload.InnerPos)
	n.innerEnd = sugar.FromSDKLex(payload.InnerEnd)
	n.body = body
	n.keyword = payload.Keyword

	return *n, nil // important: other is checking with value receiver not pointer receiver
}

var _ sugar.Node = (*SugarPlaceholderFunc)(nil)
var _ transport.NodeSerializer = (*SugarPlaceholderFunc)(nil)

type SugarPlaceholderFunc struct {
	pos      sugar.Lexeme
	end      sugar.Lexeme
	innerPos sugar.Lexeme
	innerEnd sugar.Lexeme
	body     []sugar.Lexeme
	keyword  string
}

type sugarPlaceholderFuncPayload struct {
	InnerPos sdk.Lex   `json:"innerPos"`
	InnerEnd sdk.Lex   `json:"innerEnd"`
	Body     []sdk.Lex `json:"body"`
	Keyword  string    `json:"keyword"`
}

func (n SugarPlaceholderFunc) Pos() sugar.Lexeme {
	return n.pos
}

func (n SugarPlaceholderFunc) End() sugar.Lexeme {
	return n.end
}

func (n SugarPlaceholderFunc) InnerPos() sugar.Lexeme {
	return n.innerPos
}

func (n SugarPlaceholderFunc) InnerEnd() sugar.Lexeme {
	return n.innerEnd
}

func (n SugarPlaceholderFunc) Keyword() string {
	return n.keyword
}

func (n SugarPlaceholderFunc) Body() []sugar.Lexeme {
	return n.body
}

func (n SugarPlaceholderFunc) Serialize() (*sdk.Node, error) {
	return serializeSugarPlaceholderFunc(&n)
}

const SugarPlaceholderFuncID = "lex.SugarPlaceholderFunc"

func SugarPlaceholderFuncName(keyword string) string {
	return fmt.Sprintf("__sugar_%s__", keyword)
}

func SugarPlaceholderFuncParser(keyword string) sugar.LexicalParser {
	const start, expectLParen, expectAny, running, end = "start", "expect-lparen", "expect-any", "running", "end"
	see := &sugar.LexemePredicate{}
	builder := sugar.NewNodeBuilder[SugarPlaceholderFunc]()

	const deep = "deep"
	doFail := builder.Fail
	doBegin := builder.Collect("begin", func(n *SugarPlaceholderFunc, l sugar.Lexeme) {
		builder.Error = false
		n.keyword = keyword
		n.pos = l
	})
	doCollectInnerPos := builder.Collect("inner-pos", func(n *SugarPlaceholderFunc, l sugar.Lexeme) {
		n.innerPos = l
	})
	doCollectBody := builder.Collect("body", func(n *SugarPlaceholderFunc, l sugar.Lexeme) {
		n.body = append(n.body, l)
	})
	doIncDeep := builder.Collect("inc", func(n *SugarPlaceholderFunc, l sugar.Lexeme) {
		builder.CounterInc(deep)
	})
	doDecDeep := builder.Collect("inc", func(n *SugarPlaceholderFunc, l sugar.Lexeme) {
		builder.CounterDec(deep)
		n.innerEnd = l
		if builder.Counter(deep) != 0 {
			n.body = append(n.body, l)
		}
	})
	doCollect := builder.Collect("end", func(n *SugarPlaceholderFunc, l sugar.Lexeme) {
		n.end = l
	})
	doAtStatementBoundary := func(lex sugar.Lexeme) {
		if builder.Counter(deep) == 0 {
			doCollect(lex)
		} else {
			doFail(lex)
		}
	}

	table := sugar.NewTransitionTable[string](CallExprID)
	table.
		Add(start, see.IdentMatch(SugarPlaceholderFuncName(keyword)), expectLParen, doBegin).
		Add(start, see.Any, end, doFail).
		Add(expectLParen, see.LeftParen, expectAny, doIncDeep).
		Add(expectLParen, see.Any, end, doFail).
		Add(expectAny, see.LeftParen, running, doIncDeep, doCollectInnerPos, doCollectBody).
		Add(expectAny, see.RightParen, end, doFail).
		Add(expectAny, see.StatementBoundary, end, doFail).
		Add(expectAny, see.Any, running, doCollectInnerPos, doCollectBody).
		Add(running, see.LeftParen, running, doIncDeep, doCollectBody).
		Add(running, see.RightParen, running, doDecDeep).
		Add(running, see.StatementBoundary, end, doAtStatementBoundary).
		Add(running, see.Any, running, doCollectBody)

	return sugar.NewLexicalParser(SugarPlaceholderFuncID, table, start, end, builder)
}
