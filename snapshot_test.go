package sugar

import (
	"bytes"
	"testing"
)

type mockSugar struct {
	Sugar
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

func (m *mockSugar) PrepareSource(id string, source []byte) {
}

func (m *mockSugar) CleanUp(sourceID, scopeID string) {
}

func (m *mockSugar) StructuralTransform(sourceID string, n Node) ([]byte, error) {
	return []byte(m.structureTransformed), nil
}

var _ Sugar = (*mockSugar)(nil)

func Test_Snapshot_doStructuralTransform(t *testing.T) {
	source := []byte("hello SUGAR world")
	//                0123456789...
	// SUGAR is at offset 6, len 5, ends at 11

	s := &mockSugar{
		pos:                  Lexeme{Offset: 6},
		end:                  Lexeme{Offset: 11},
		structureTransformed: "REPLACED",
	}

	ss := &Snapshot{
		source: source,
		nodes: []parsedNode{
			{
				sugar: s,
				node:  s,
			},
		},
	}
	_ = ss.StructuralTransform()

	got := ss.t1output
	want := []byte("hello REPLACED world")

	if !bytes.Equal(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}
