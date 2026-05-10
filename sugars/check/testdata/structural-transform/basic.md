## Structure transformation

given the input is:

```gos
// file: input.gos
package example

import (
	"fmt"
	"strconv"
)

func test() {
	check doSomething()
	
	x := check doSomething()
	y := check strconv.Atoi("123")

	fmt.Println(x, y)
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

func test() {
	__sugar_check__(doSomething())
	
	x := __sugar_check__(doSomething())
	y := __sugar_check__(strconv.Atoi("123"))

	fmt.Println(x, y)
}
```
