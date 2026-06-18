package startcmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"golang.org/x/exp/jsonrpc2"
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/internal/sdk"
	"nhatp.com/go/sugar/internal/sdk/transport"
)

type Handler struct {
	sugar sugar.Sugar
	log   *slog.Logger
}

func (h *Handler) Handle(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	switch req.Method {
	case sdk.Method.ID:
		return h.sugar.ID(), nil

	case sdk.Method.Binary:
		return h.sugar.Binary(), nil

	case sdk.Method.Parse:
		source, output, err := parseSource(req, h.sugar.Parse)
		if err != nil {
			return nil, err
		}
		h.log.Debug("parse", slog.String("source", source), slog.Any("output", output))
		return output, nil

	case sdk.Method.Restore:
		source, output, err := parseSource(req, h.sugar.Restore)
		if err != nil {
			return nil, err
		}
		h.log.Debug("restore", slog.String("source", source), slog.Any("output", output))
		return output, nil

	case sdk.Method.PrepareSource:
		var params sdk.PrepareSourceParam
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return nil, err
		}
		h.sugar.PrepareSource(params.ID, []byte(params.Source))
		return true, nil

	case sdk.Method.PrepareSemanticScope:
		var params sdk.PrepareSemanticScopeParam[sugar.SemanticScope]
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return nil, err
		}
		h.sugar.PrepareSemanticScope(params.ID, params.SemanticScope)
		return true, nil

	case sdk.Method.CleanUp:
		var params sdk.CleanUpParam
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return nil, err
		}
		h.sugar.CleanUp(params.SourceID, params.ScopeID)
		return true, nil

	case sdk.Method.StructuralTransform:
		out, err := handleTransform(req, func(param sdk.TransformParam, node sugar.Node) ([]byte, error) {
			return h.sugar.StructuralTransform(param.SourceID, node)
		})
		if err != nil {
			return nil, err
		}
		return string(out), nil

	case sdk.Method.SemanticTransformer:
		out, err := handleTransform(req, func(param sdk.TransformParam, node sugar.Node) ([]byte, error) {
			return h.sugar.SemanticTransformer(param.SourceID, param.ScopeID, node)
		})
		if err != nil {
			return nil, err
		}
		return string(out), nil

	case sdk.Method.RestoreTransform:
		out, err := handleTransform(req, func(param sdk.TransformParam, node sugar.Node) ([]byte, error) {
			return h.sugar.RestoreTransform(param.SourceID, node)
		})
		if err != nil {
			return nil, err
		}
		return string(out), nil

	default:
		return nil, jsonrpc2.ErrMethodNotFound
	}
}

func handleTransform(req *jsonrpc2.Request, fn func(sdk.TransformParam, sugar.Node) ([]byte, error)) ([]byte, error) {
	var params sdk.TransformParam
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, err
	}

	node, err := fromSDKNode(params.Node)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, fmt.Errorf("invalid node after deserialization")
	}
	return fn(params, node)
}

func parseSource(req *jsonrpc2.Request, fn func([]byte) []sugar.Node) (string, []*sdk.Node, error) {
	var source string
	err := json.Unmarshal(req.Params, &source)
	if err != nil {
		return "", nil, err
	}

	nodes := fn([]byte(source))
	var output []*sdk.Node
	for _, n := range nodes {
		v, err := toSDKNode(n)
		if err != nil {
			return source, nil, err
		}
		output = append(output, v)
	}
	return source, output, nil
}

func toSDKNode(in sugar.Node) (*sdk.Node, error) {
	s, ok := in.(transport.NodeSerializer)
	if ok {
		return s.Serialize()
	}

	raw, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	pos := in.Pos()
	end := in.End()
	v := sdk.Node{
		Pos:     pos.ToSDKLex(),
		End:     end.ToSDKLex(),
		Payload: raw,
	}
	return &v, nil
}

func fromSDKNode(in sdk.Node) (sugar.Node, error) {
	return transport.DeserializeNode(in)
}
