## Structural Transform Golden Test

Write something for all the whole test suite, this is shared between all test cases

Phase One: Structural Transform 

- Input: .gos file
- Output: valid Go that preserves the structure of the sugar
- Good enough for `gofmt` and `go/parser` / `packages.Load()` to process
- Semantic correctness is not required
- Builds a lightweight source map (offsets only, transient, never written to disk)

### Test case A

> Input is define by a codeblock with `file: input.gos`

Given the input is:

```go
// file: input.gos
package example

import (
	"fmt"
	"strconv"
)

func test() {
	x := check doSomething()
	y := check strconv.Atoi("123")

	fmt.Println(x, y)
}
```

> Golden output is defined by a codeblock with `// golden-file: output.go`

the Structural Transform output is:

```go
// golden-file: output.go
package example

import (
	"fmt"
	"strconv"
)

func test() {
	x := __sugar_check__(doSomething())
	y := __sugar_check__(strconv.Atoi("123"))

	fmt.Println(x, y)
}
```

### Test case B

Same structure as test case A, you can have multiple cases in one test suite.
