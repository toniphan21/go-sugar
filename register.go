package sugar

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"slices"

	"nhatp.com/go/sugar/internal/sdk"
)

// --- registration

var plugins map[string]Sugar
var processes map[string]sdk.Process

func Register(plugin Sugar) {
	if plugin == nil {
		return
	}

	if plugins == nil {
		plugins = make(map[string]Sugar)
	}
	plugins[plugin.ID()] = plugin
}

func RegisterBinary(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	if processes == nil {
		processes = make(map[string]sdk.Process)
	}

	process := sdk.NewProcess(path, "")
	processes[path] = process
	binary := &binaryClient{
		process: process,
	}

	id := binary.ID()
	if binary.err != nil {
		return binary.err
	}

	if plugins == nil {
		plugins = make(map[string]Sugar)
	}
	plugins[id] = binary
	return nil
}

func StopBinaryProcesses() {
	for _, v := range processes {
		_ = v.Close()
	}
}

func registered() []Sugar {
	return slices.Collect(maps.Values(plugins))
}

// --- sugar binary transport

type binaryNodeWrapper struct {
	node sdk.Node
}

func (w *binaryNodeWrapper) Pos() Lexeme {
	return FromSDKLex(w.node.Pos)
}

func (w *binaryNodeWrapper) End() Lexeme {
	return FromSDKLex(w.node.End)
}

type binaryClient struct {
	process sdk.Process
	err     error
}

func (b *binaryClient) ID() string {
	b.err = nil
	var result string

	b.err = b.process.Call(context.Background(), sdk.Method.ID, nil, &result)
	return result
}

func (b *binaryClient) Binary() BinaryInfo {
	b.err = nil
	var result BinaryInfo
	b.err = b.process.Call(context.Background(), sdk.Method.Binary, nil, &result)

	return result
}

func (b *binaryClient) doParseNode(source []byte, rpc string) []Node {
	b.err = nil
	s, err := json.Marshal(string(source))
	if err != nil {
		b.err = err
		return nil
	}

	var nodes []sdk.Node
	b.err = b.process.Call(context.Background(), rpc, s, &nodes)
	if b.err != nil {
		fmt.Println(b.err.Error())
	}

	var result []Node
	for _, v := range nodes {
		result = append(result, &binaryNodeWrapper{node: v})
	}
	return result
}

func (b *binaryClient) Parse(source []byte) []Node {
	return b.doParseNode(source, sdk.Method.Parse)
}

func (b *binaryClient) Restore(source []byte) []Node {
	return b.doParseNode(source, sdk.Method.Restore)
}

func (b *binaryClient) PrepareSource(id string, source []byte) {
	b.err = nil
	param, err := json.Marshal(sdk.PrepareSourceParam{
		ID:     id,
		Source: string(source),
	})
	if err != nil {
		b.err = err
		return
	}
	b.err = b.process.Call(context.Background(), sdk.Method.PrepareSource, param, nil)
}

func (b *binaryClient) PrepareSemanticScope(id string, scope SemanticScope) {
	b.err = nil
	param, err := json.Marshal(sdk.PrepareSemanticScopeParam[SemanticScope]{
		ID:            id,
		SemanticScope: scope,
	})
	if err != nil {
		b.err = err
		return
	}
	b.err = b.process.Call(context.Background(), sdk.Method.PrepareSemanticScope, param, nil)
}

func (b *binaryClient) CleanUp(sourceID string, scopeID string) {
	b.err = nil
	param, err := json.Marshal(sdk.CleanUpParam{
		SourceID: sourceID,
		ScopeID:  scopeID,
	})
	if err != nil {
		b.err = err
		return
	}
	b.err = b.process.Call(context.Background(), sdk.Method.CleanUp, param, nil)
}

func (b *binaryClient) doTransform(rpc string, n Node, fn func(sdk.Node) sdk.TransformParam) ([]byte, error) {
	b.err = nil
	bn, ok := n.(*binaryNodeWrapper)
	if !ok {
		return nil, fmt.Errorf("expected binaryNodeWrapper")
	}
	param, err := json.Marshal(fn(bn.node))
	if err != nil {
		b.err = err
		return nil, err
	}

	var output string
	err = b.process.Call(context.Background(), rpc, param, &output)

	return []byte(output), err
}

func (b *binaryClient) StructuralTransform(sourceID string, n Node) ([]byte, error) {
	return b.doTransform(sdk.Method.StructuralTransform, n, func(node sdk.Node) sdk.TransformParam {
		return sdk.TransformParam{
			SourceID: sourceID,
			Node:     node,
		}
	})
}

func (b *binaryClient) SemanticTransformer(sourceID string, scopeID string, n Node) ([]byte, error) {
	return b.doTransform(sdk.Method.SemanticTransformer, n, func(node sdk.Node) sdk.TransformParam {
		return sdk.TransformParam{
			SourceID: sourceID,
			ScopeID:  scopeID,
			Node:     node,
		}
	})
}

func (b *binaryClient) RestoreTransform(sourceID string, n Node) ([]byte, error) {
	return b.doTransform(sdk.Method.RestoreTransform, n, func(node sdk.Node) sdk.TransformParam {
		return sdk.TransformParam{
			SourceID: sourceID,
			Node:     node,
		}
	})
}

var _ Sugar = (*binaryClient)(nil)
