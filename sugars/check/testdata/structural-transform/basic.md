## Structure transformation

given a go module

```go.mod
module github.com/you/repo

go 1.24
```

given the input is:

```gos
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

the Structural Transform output is:

```go
// golden-file: output.go
package example

import (
	"fmt"
	"strconv"
)

func test() error {
	__sugar_check__(doSomething())
	
	x := __sugar_check__(strconv.Atoi("123"))

	fmt.Println(x)
	return nil
}

func doSomething() error {
	return nil
}
```
