package require

import (
	"encoding/json"
	"go/ast"
	"go/types"

	"github.com/oklog/ulid/v2"
	"golang.org/x/tools/go/packages"
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/internal/sdk"
	"nhatp.com/go/sugar/internal/sdk/transport"
	"nhatp.com/go/sugar/sa"
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

	enclosingParams []*types.Var
	callReturns     []types.Type
}

func (n *node) Pos() sugar.Lexeme {
	return n.pos
}

func (n *node) End() sugar.Lexeme {
	return n.end
}

func (n *node) Serialize() (*sdk.Node, error) {
	return serializeNode(n)
}

func (n *node) SemanticAnalysis(pkg *packages.Package, smap *sugar.SourceMap) error {
	sa.InspectPlaceholderFunc(pkg, smap, n, keyword, func(file *ast.File, call *ast.CallExpr) bool {
		// reset semantic data
		n.callReturns = nil
		n.enclosingParams = nil

		// the inner expression is the first arg of __sugar_require__
		innerExpr := call.Args[0]
		tv, ok := pkg.TypesInfo.Types[innerExpr]
		if ok {
			if tuple, ok := tv.Type.(*types.Tuple); ok {
				for i := 0; i < tuple.Len(); i++ {
					n.callReturns = append(n.callReturns, tuple.At(i).Type())
				}
			}
		}

		// find the enclosing function
		enclosing := sa.FindEnclosingFunc(file, call.Pos())
		if enclosing == nil {
			return true
		}
		fn := pkg.TypesInfo.Defs[enclosing.Name].(*types.Func)
		sig, ok := fn.Type().(*types.Signature)
		if !ok {
			return true
		}
		if sig.Params() != nil {
			params := sig.Params()
			for i := 0; i < params.Len(); i++ {
				n.enclosingParams = append(n.enclosingParams, params.At(i))
			}
		}
		return true
	})

	return nil
}

func (n *node) findTestingParam() (string, bool) {
	for _, v := range n.enclosingParams {
		t := v.Type()
		if sa.IsTestingTB(t) || sa.IsTestingReceiver(t) {
			return v.Name(), true
		}
	}
	return "", false
}

func (n *node) scanMessageVerbs(s string) (hasS, hasV bool) {
	for i := 0; i < len(s); i++ {
		if s[i] != '%' || i+1 >= len(s) {
			continue
		}
		switch s[i+1] {
		case '%':
			i++ // literal percent, skip
		case 's':
			hasS = true
			i++
		case 'v':
			hasV = true
			i++
		}
	}
	return
}

var _ sugar.SemanticNode = (*node)(nil)
var _ transport.NodeSerializer = (*node)(nil)
