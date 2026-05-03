package sugar

import (
	"fmt"
	"testing"
)

func Test_Lex(t *testing.T) {
	source := `func test() {
	var a int
	x := check strconv.Atoi(
		a,
		b,
	)
}`
	lex := Lex([]byte(source))
	for _, v := range lex {
		fmt.Printf("%d|%v|%v\n", v.Offset, v.Tok, v.Lit)
	}
}
