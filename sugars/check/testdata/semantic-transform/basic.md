## Semantic Transform

given a go module

```go.mod
module github.com/you/repo

go 1.24
```

### Without zero value in enclosing function

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

the Transform output is

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
