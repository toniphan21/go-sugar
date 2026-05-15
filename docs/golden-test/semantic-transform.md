## Semantic Transform Golden Test

Semantic transform (T2) converts `.gos` sugar syntax into compilable Go using type information gathered from T1 output
via `packages.Load()`. The goal of this test is not to produce correctly formatted or indented Go code - it is to
produce output that is valid enough to feed to `gopls` and proxy diagnostics back to the editor. Formatting is handled
separately by `gofmt` or `goimports`.

Note: emitted code from sugar expansion always starts at indent level zero, regardless of the surrounding context. This
is intentional - formatting is not T2's concern.

All content under this section is shared across all test cases below. Most commonly, this is the `go.mod` file that
defines the module boundary.

```go.mod
module github.com/you/repo

go 1.24
```

### Test case A

> Input is defined by a code block with `// file: input.gos`

Given the input:

```go
// file: input.gos
package example

import (
	"fmt"
	"strconv"
)

func test() error {
	check doSomething()

	x := check strconv.Atoi("123")

	fmt.Println(x)
	return nil
}

func doSomething() error {
	return nil
}
```

> Golden output is defined by a code block with `// golden-file: output.go`

The semantic transform output is:

```go
// golden-file: output.go
package example

import (
	"fmt"
	"strconv"
)

func test() error {
	err := doSomething()
if err != nil {
	return err
}

	x, err := strconv.Atoi("123")
if err != nil {
	return err
}

	fmt.Println(x)
	return nil
}

func doSomething() error {
	return nil
}
```

### Test case B

> This test case follows the same structure as test case A.