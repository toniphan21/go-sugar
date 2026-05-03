package check

import (
	"testing"
)

func Test_Scan_Join(t *testing.T) {
	source := `x := check strconv.Atoi(a)`
	tokens := Scan([]byte(source))
	out := Join(tokens)

	if out != source {
		t.Fatalf("\nexpected: %s\nactual: %s", source, out)
	}
}
