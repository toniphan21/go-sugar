package sugar

import (
	"bytes"
	"testing"

	"golang.org/x/tools/go/packages"
)

type mockSugar struct {
	pos                  Lexeme
	end                  Lexeme
	structureTransformed string
}

func (m *mockSugar) Pos() Lexeme {
	return m.pos
}

func (m *mockSugar) End() Lexeme {
	return m.end
}

func (m *mockSugar) StructuralTransform(source []byte, lexemes []Lexeme) []byte {
	return []byte(m.structureTransformed)
}

func (m *mockSugar) SemanticAnalysis(pkg *packages.Package, smap *SourceMap) error {
	return nil
}

func (m *mockSugar) SemanticTransformer(source []byte, lexemes []Lexeme) []byte {
	return nil
}

var _ Sugar = (*mockSugar)(nil)

func Test_Snapshot_doStructuralTransform(t *testing.T) {
	source := []byte("hello SUGAR world")
	//                0123456789...
	// SUGAR is at offset 6, len 5, ends at 11

	sugars := []Sugar{
		&mockSugar{
			pos:                  Lexeme{Offset: 6},
			end:                  Lexeme{Offset: 11},
			structureTransformed: "REPLACED",
		},
	}

	ss := &Snapshot{
		source: source,
	}
	ss.doStructuralTransform(sugars)

	got := ss.t1output
	want := []byte("hello REPLACED world")

	if !bytes.Equal(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}
