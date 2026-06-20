package check

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
	CheckPos    sdk.Lex  `json:"checkPos"`
	CheckEnd    sdk.Lex  `json:"checkEnd"`
	Identifiers []string `json:"identifiers"`
}

func serializeNode(n *node) (*sdk.Node, error) {
	payload, err := json.Marshal(&nodePayload{
		CheckPos:    n.checkPos.ToSDKLex(),
		CheckEnd:    n.checkEnd.ToSDKLex(),
		Identifiers: n.identifiers,
	})
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

	n.checkPos = sugar.FromSDKLex(payload.CheckPos)
	n.checkEnd = sugar.FromSDKLex(payload.CheckEnd)
	n.identifiers = payload.Identifiers
	return n, nil
}

var _ sugar.SemanticNode = (*node)(nil)
var _ transport.NodeSerializer = (*node)(nil)

type node struct {
	pos              sugar.Lexeme
	end              sugar.Lexeme
	checkPos         sugar.Lexeme
	checkEnd         sugar.Lexeme
	identifiers      []string
	enclosingReturns []types.Type
	callReturns      []types.Type
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
		n.enclosingReturns = nil
		n.callReturns = nil

		// find the enclosing function
		enclosing := sa.FindEnclosingFunc(file, call.Pos())
		if enclosing == nil {
			return true
		}
		if enclosing.Type.Results != nil {
			for _, field := range enclosing.Type.Results.List {
				t := pkg.TypesInfo.TypeOf(field.Type)
				n.enclosingReturns = append(n.enclosingReturns, t)
			}
		}

		// the inner expression is the first arg of __sugar_check__
		innerExpr := call.Args[0]
		tv, ok := pkg.TypesInfo.Types[innerExpr]
		if ok {
			if tuple, ok := tv.Type.(*types.Tuple); ok {
				for i := 0; i < tuple.Len(); i++ {
					n.callReturns = append(n.callReturns, tuple.At(i).Type())
				}
			}
		}
		return true
	})
	return nil
}
