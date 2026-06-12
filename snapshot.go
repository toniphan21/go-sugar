package sugar

import (
	"bytes"
	"cmp"
	"crypto/sha256"
	"slices"
)

func newSnapshot(content []byte) *Snapshot {
	hash := sha256.Sum256(content)

	return &Snapshot{hash: hash, source: content}
}

type Snapshot struct {
	hash     [32]byte
	source   []byte
	nodes    []parsedNode
	t1output []byte
	t1smap   *SourceMap
	t2output []byte
	t2smap   *SourceMap
}

func (s *Snapshot) doParse() []parsedNode {
	if s.nodes == nil {
		for _, v := range registered() {
			nodes := v.Parse(s.source)
			for _, n := range nodes {
				s.nodes = append(s.nodes, parsedNode{
					sugar: v,
					node:  n,
				})
			}
		}

		slices.SortFunc(s.nodes, func(a, b parsedNode) int {
			return cmp.Compare(a.node.Pos().Offset, b.node.Pos().Offset)
		})
	}
	return s.nodes
}

func (s *Snapshot) doTransform(pn []parsedNode, fn func(parsedNode) ([]byte, error)) ([]byte, *SourceMap) {
	out := bytes.Buffer{}
	smap := &SourceMap{}
	cursor := 0
	for _, v := range pn {
		out.Write(s.source[cursor:v.node.Pos().Offset])

		goStart := out.Len()
		if transformed, err := fn(v); err != nil {
			out.Write(s.source[v.node.Pos().Offset:v.node.End().Offset]) // pass-through
		} else {
			out.Write(transformed)
		}
		goEnd := out.Len()

		smap.Entries = append(smap.Entries, Entry{
			Sugar: Region{
				Pos: Position{Offset: v.node.Pos().Offset},
				End: Position{Offset: v.node.End().Offset},
			},
			Go: Region{
				Pos: Position{Offset: goStart},
				End: Position{Offset: goEnd},
			},
			Kind: KindExpand,
		})

		cursor = v.node.End().Offset
	}

	out.Write(s.source[cursor:])

	return out.Bytes(), smap
}

func (s *Snapshot) StructuralTransform() []byte {
	if s.t1smap == nil {
		parsedNodes := s.doParse()

		sourceId := makeSourceID(s.source)
		for _, v := range parsedNodes {
			v.sugar.PrepareSource(sourceId, s.source)
		}

		out, smap := s.doTransform(parsedNodes, func(pn parsedNode) ([]byte, error) {
			return pn.sugar.StructuralTransform(sourceId, pn.node)
		})

		for _, v := range parsedNodes {
			v.sugar.CleanUp(sourceId, "")
		}

		smap.buildByGo()

		s.t1output = out
		s.t1smap = smap
	}
	return s.t1output
}

func (s *Snapshot) SemanticTransform(module ModuleScope, file FileScope) ([]byte, error) {
	if s.t2smap == nil {
		parsedNodes := s.doParse()

		sourceId := makeSourceID(s.source)
		scopeId := makeScopeID(module, file)
		for _, v := range parsedNodes {
			v.sugar.PrepareSource(sourceId, s.source)
			v.sugar.PrepareSemanticScope(scopeId, SemanticScope{
				ModuleScope: module,
				FileScope:   file,
			})
		}

		out, smap := s.doTransform(parsedNodes, func(pn parsedNode) ([]byte, error) {
			return pn.sugar.SemanticTransformer(sourceId, scopeId, pn.node)
		})

		for _, v := range parsedNodes {
			v.sugar.CleanUp(sourceId, scopeId)
		}

		smap.buildBySugar()
		smap.buildByGo()

		s.t2output = out
		s.t2smap = smap
	}
	return s.t2output, nil
}

func (s *Snapshot) Hash() [32]byte {
	return s.hash
}

func (s *Snapshot) SugarToGo(line, column int) (int, int, error) {
	return line, column, nil
}

func (s *Snapshot) GoToSugar(line, column int) (int, int, error) {
	return line, column, nil
}

var _ snapshotAPI = (*Snapshot)(nil)
