package require

import (
	"encoding/json"

	"github.com/oklog/ulid/v2"
	"golang.org/x/tools/go/packages"
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/internal/sdk"
	"nhatp.com/go/sugar/internal/sdk/transport"
)

func init() {
	transport.RegisterNodeDeserializer(NodeType, deserializeNode)
}

type nodePayload struct {
	RequirePos  sdk.Lex  `json:"checkPos"`
	RequireEnd  sdk.Lex  `json:"checkEnd"`
	MessagePos  *sdk.Lex `json:"messagePos,omitempty"`
	MessageEnd  *sdk.Lex `json:"messageEnd,omitempty"`
	Message     *string  `json:"message,omitempty"`
	Identifiers []string `json:"identifiers"`
}

func serializeNode(n *node) (*sdk.Node, error) {
	np := &nodePayload{
		RequirePos:  n.requirePos.ToSDKLex(),
		RequireEnd:  n.requireEnd.ToSDKLex(),
		Message:     n.message,
		Identifiers: n.identifiers,
	}
	if n.messagePos != nil {
		np.MessagePos = new(n.messagePos.ToSDKLex())
	}
	if n.messageEnd != nil {
		np.MessageEnd = new(n.messageEnd.ToSDKLex())
	}

	payload, err := json.Marshal(np)
	if err != nil {
		return nil, err
	}

	return &sdk.Node{
		ID:      ulid.Make().String(),
		Type:    NodeType,
		Pos:     n.pos.ToSDKLex(),
		End:     n.end.ToSDKLex(),
		Payload: payload,
	}, nil
}

func deserializeNode(in sdk.Node) (sugar.Node, error) {
	n := &node{}
	n.pos = sugar.FromSDKLex(in.Pos)
	n.end = sugar.FromSDKLex(in.End)

	var payload nodePayload
	if err := json.Unmarshal(in.Payload, &payload); err != nil {
		return nil, err
	}

	n.requirePos = sugar.FromSDKLex(payload.RequirePos)
	n.requireEnd = sugar.FromSDKLex(payload.RequireEnd)
	n.message = payload.Message
	n.identifiers = payload.Identifiers
	if payload.MessagePos != nil {
		n.messagePos = new(sugar.FromSDKLex(*payload.MessagePos))
	}
	if payload.MessageEnd != nil {
		n.messageEnd = new(sugar.FromSDKLex(*payload.MessageEnd))
	}
	return n, nil
}

type node struct {
	pos         sugar.Lexeme
	end         sugar.Lexeme
	requirePos  sugar.Lexeme
	requireEnd  sugar.Lexeme
	messagePos  *sugar.Lexeme
	messageEnd  *sugar.Lexeme
	message     *string
	identifiers []string
}

func (n *node) Pos() sugar.Lexeme {
	return n.pos
}

func (n *node) End() sugar.Lexeme {
	return n.end
}

func (n *node) SemanticAnalysis(pkg *packages.Package, smap *sugar.SourceMap) error {
	return nil
}

func (n *node) Serialize() (*sdk.Node, error) {
	return serializeNode(n)
}

var _ sugar.SemanticNode = (*node)(nil)
var _ transport.NodeSerializer = (*node)(nil)
