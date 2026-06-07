## Format Pipeline

given a go module
```go.mod
module github.com/you/repo

go 1.24
```

given the input:
```gos
// file: input.gos
package example

func test() {x := require doSomething()
	y := require strconv.Atoi(            "123")     "it is not a number"
}
```

the formatted output is:
```gos
// golden-file: output.gos
package example

func test() {
	x := require doSomething()
	y := require strconv.Atoi("123") "it is not a number"
}
```
