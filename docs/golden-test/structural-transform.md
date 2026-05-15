## Structural Transform Golden Test

Structural transform (T1) converts `.gos` sugar syntax into valid Go with placeholders. The output is good enough for
`gofmt` and `go/parser` / `packages.Load()` to process. Semantic correctness is not required. T1 also builds a
lightweight source map (offsets only, transient, never written to disk).

All content under this section is shared across all test cases below. Most commonly, this is the `go.mod` file that
defines the module boundary.

```go.mod
module github.com/you/repo

go 1.24
```

### Test case A

> Input is defined by a code block with `// file: input.gos`

Given the input:

```gos
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

> Golden output is defined by a code block with `// golden-file: output.go`

The structural transform output is:

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

Same structure as test case A. Multiple cases can appear in one test suite.