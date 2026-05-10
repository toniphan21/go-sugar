package lextest

import (
	"fmt"
	"testing"

	"nhatp.com/go/sugar"
)

func Nodes[T any](items ...T) []T {
	return items
}

func AssertNodes[T, E any](
	t testing.TB,
	code string,
	expected []T, actual []E,
	compare func(T, E) (string, bool),
) {
	t.Helper()

	if len(actual) != len(expected) {
		t.Fatalf("len(result) = %d, want %d\n%v", len(actual), len(expected), msgForLexViewer(code))
	}

	for i, v := range actual {
		msg, valid := compare(expected[i], v)
		if !valid {
			t.Fatalf("index %d: %v\n%v", i, msg, msgForLexViewer(code))
		}
	}
}

func msgForLexViewer(code string) string {
	return fmt.Sprintf(
		"Tip: copy line below to tools/lexeme-viewer at http://localhost:%d for debugging\n%v\n",
		sugar.ToolLexemeViewerDefaultPort,
		FormatCodeForLexViewer(code),
	)
}

func CompareOptionalPos(name string, want int, got *sugar.Lexeme) (string, bool) {
	gp := 0
	if got != nil {
		gp = int(got.Pos)
	}
	if want != gp {
		return fmt.Sprintf("got %v=%d, want %d", name, gp, want), false
	}
	return "", true
}
