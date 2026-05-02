package check

import (
	"testing"
)

func Test_Scan_Join(t *testing.T) {
	source := `x := check strconv.Atoi(a)`
	tokens := Scan([]byte(source))
	out := Join(tokens)

	expected := "x :=       strconv.Atoi(a)"
	if out != expected {
		t.Fatalf("\nexpected: %s\nactual: %s", expected, out)
	}
}
