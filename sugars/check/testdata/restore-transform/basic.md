## Restore Transform

given a go module

```go.mod
module github.com/you/repo

go 1.24
```

given the input
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

the Restore Transform output is:

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
