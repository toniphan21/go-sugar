package sugar

import (
	"bytes"
	"cmp"
	"crypto/sha256"
	"slices"

	"golang.org/x/tools/go/packages"
)

func newSnapshot(content []byte) *Snapshot {
	hash := sha256.Sum256(content)

	return &Snapshot{hash: hash, source: content}
}

type Snapshot struct {
	hash     [32]byte
	source   []byte
	lexemes  []Lexeme
	sugars   []Sugar
	t1output []byte
	t1smap   *SourceMap
	t2output []byte
	t2smap   *SourceMap
}

func (s *Snapshot) scan() []Lexeme {
	if s.lexemes == nil {
		s.lexemes = Lex(s.source)
	}
	return s.lexemes
}

func (s *Snapshot) doParseSugars() []Sugar {
	if s.sugars == nil {
		s.sugars = parseSugars(s.scan())
	}
	return s.sugars
}

func (s *Snapshot) doTransform(sugars []Sugar, fn func(Sugar, []byte, []Lexeme) []byte) ([]byte, *SourceMap) {
	out := bytes.Buffer{}
	smap := &SourceMap{}
	cursor := 0

	for _, v := range sugars {
		out.Write(s.source[cursor:v.Pos().Offset])

		goStart := out.Len()
		transformed := fn(v, s.source, s.scan())
		out.Write(transformed)
		goEnd := out.Len()

		smap.Entries = append(smap.Entries, Entry{
			Sugar: Region{
				Pos: Position{Offset: v.Pos().Offset},
				End: Position{Offset: v.End().Offset},
			},
			Go: Region{
				Pos: Position{Offset: goStart},
				End: Position{Offset: goEnd},
			},
			Kind: KindExpand,
		})

		cursor = v.End().Offset
	}

	out.Write(s.source[cursor:])

	return out.Bytes(), smap
}

func (s *Snapshot) doStructuralTransform(sugars []Sugar) {
	out, smap := s.doTransform(sugars, func(v Sugar, s []byte, l []Lexeme) []byte {
		return v.StructuralTransform(s, l)
	})

	smap.buildByGo()

	s.t1output = out
	s.t1smap = smap
}

func (s *Snapshot) doSemanticTransform(sugars []Sugar) {
	out, smap := s.doTransform(sugars, func(v Sugar, s []byte, l []Lexeme) []byte {
		return v.SemanticTransformer(s, l)
	})

	smap.buildBySugar()
	smap.buildByGo()

	s.t2output = out
	s.t2smap = smap
}

func (s *Snapshot) structuralTransform() error {
	if s.t1smap != nil {
		return nil
	}

	sugars := s.doParseSugars()
	slices.SortFunc(sugars, func(a, b Sugar) int {
		return cmp.Compare(a.Pos().Offset, b.Pos().Offset)
	})

	s.doStructuralTransform(sugars)
	return nil
}

func (s *Snapshot) semanticAnalysis(pkg *packages.Package) error {
	sugars := s.doParseSugars()
	for _, v := range sugars {
		if err := v.SemanticAnalysis(pkg, s.t1smap); err != nil {
			return err
		}
	}
	return nil
}

func (s *Snapshot) semanticTransform() error {
	s.doSemanticTransform(s.doParseSugars())

	return nil
}

func (s *Snapshot) StructuralTransform() []byte {
	return s.t1output
}

func (s *Snapshot) SemanticTransform() []byte {
	return s.t2output
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
