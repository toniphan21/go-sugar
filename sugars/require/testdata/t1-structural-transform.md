## Structural Transform

given a go module
```go.mod
module github.com/you/repo

go 1.24
```

### Basic

#### Function call

given the input is
```gos
// file: input.gos
package example

func test() {
	require setUp()
	svc := require newService()
}
```

the T1 - Structural Transform output is
```go
// golden-file: output.go
package example

func test() {
	__sugar_require__(setUp())
	svc := __sugar_require__(newService())
}
```

#### Function call with message

given the input is
```gos
// file: input.gos
package example

func test() {
	require setUp() "cannot setup"
	svc := require newService() "cannot make service"
}
```

the T1 - Structural Transform output is
```go
// golden-file: output.go
package example

func test() {
	__sugar_require__(setUp() ,"cannot setup")
	svc := __sugar_require__(newService() ,"cannot make service")
}
```

#### Function call with multi-line message

given the input is
```gos
// file: input.gos
package example

func test() {
	require setUp() `cannot setup
please double check`
	svc := require newService() `cannot make service
please double check`
}
```

the T1 - Structural Transform output is
```go
// golden-file: output.go
package example

func test() {
	__sugar_require__(setUp() ,`cannot setup
please double check`)
	svc := __sugar_require__(newService() ,`cannot make service
please double check`)
}
```

#### Function call with params

given the input is
```gos
// file: input.gos
package example

func test() {
	require setUp(1, "any")
	svc := require newService("server")
}
```

the T1 - Structural Transform output is
```go
// golden-file: output.go
package example

func test() {
	__sugar_require__(setUp(1, "any"))
	svc := __sugar_require__(newService("server"))
}
```

#### Function call with param and message

given the input is
```gos
// file: input.gos
package example

func test() {
	require setUp("any") "cannot setup"
	svc := require newService(123) "cannot make service"
}
```

the T1 - Structural Transform output is
```go
// golden-file: output.go
package example

func test() {
	__sugar_require__(setUp("any") ,"cannot setup")
	svc := __sugar_require__(newService(123) ,"cannot make service")
}
```

#### Function call with params and multi-line message

given the input is
```gos
// file: input.gos
package example

func test() {
	require setUp(whatever) `cannot setup
please double check`
	svc := require newService(true) `cannot make service
please double check`
}
```

the T1 - Structural Transform output is
```go
// golden-file: output.go
package example

func test() {
	__sugar_require__(setUp(whatever) ,`cannot setup
please double check`)
	svc := __sugar_require__(newService(true) ,`cannot make service
please double check`)
}
```
