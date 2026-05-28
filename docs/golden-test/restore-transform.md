## Restore Transform Golden Test

Restore transform (T3) converts `.go` back to sugar.

All content under this section is shared across all test cases below. Most commonly, this is the `go.mod` file that
defines the module boundary.

```go.mod
module github.com/you/repo

go 1.24
```

### Test case A

> Input is defined by a code block with `// file: input.go`

Given the input:

```go
// file: input.go
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

> Golden output is defined by a code block with `// golden-file: output.gos`

The restore transform output is:

```gos
// golden-file: output.gos
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

### Test case B

Same structure as test case A. Multiple cases can appear in one test suite.