package sugar

import (
	"bytes"
	"cmp"
	"slices"
)

func parseRestoreNodes(source []byte) []parsedNode {
	var result []parsedNode
	for _, v := range registered() {
		nodes := v.Restore(source)
		for _, n := range nodes {
			result = append(result, parsedNode{
				sugar: v,
				node:  n,
			})
		}
	}

	slices.SortFunc(result, func(a, b parsedNode) int {
		return cmp.Compare(a.node.Pos().Offset, b.node.Pos().Offset)
	})
	return result
}

func RestoreTransform(source []byte) []byte {
	parsedNodes := parseRestoreNodes(source)

	sourceId := makeSourceID(source)
	for _, v := range parsedNodes {
		v.sugar.PrepareSource(sourceId, source)
	}

	out := bytes.Buffer{}
	cursor := 0
	for _, v := range parsedNodes {
		out.Write(source[cursor:v.node.Pos().Offset])
		if transformed, err := v.sugar.RestoreTransform(sourceId, v.node); err != nil {
			out.Write(source[v.node.Pos().Offset:v.node.End().Offset]) // pass-through
		} else {
			out.Write(transformed)
		}
		cursor = v.node.End().Offset
	}

	out.Write(source[cursor:])

	return out.Bytes()
}
