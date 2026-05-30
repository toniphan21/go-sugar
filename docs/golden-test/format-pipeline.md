## Format Pipeline Golden Test

Format Pipeline format `.gos` using process T1 -> gofmt -> T3.

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

func test() {x := check doSomething()
	y := check strconv.Atoi(            "123")

	fmt.Println(x, y)
}
```

> Golden output is defined by a code block with `// golden-file: output.gos`

The formatted output is:

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