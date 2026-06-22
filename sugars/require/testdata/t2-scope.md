## Semantic Transform - Scope

given a go module

```go.mod
module github.com/you/repo

go 1.24
```

### prod

given the input

```gos
// file: input.gos
package example

func setUp() error { return nil }

func DoSomething() {
	require setUp()
}
```

the T2 output is

```
// golden-file: output.go
package example

func setUp() error { return nil }

func DoSomething() {
	if err := setUp(); err != nil {
	panic("setUp(): " + err.Error())
}
}
```

### test

given the input

```gos
// file: input.gos
package example

import "testing"

func setUp() error { return nil }

func Test(t *testing.T) {
	require setUp()
}
```

the T2 output is

```
// golden-file: output.go
package example

import "testing"

func setUp() error { return nil }

func Test(t *testing.T) {
	if err := setUp(); err != nil {
	t.Fatalf("%s: %v", "setUp()", err)
}
}
```

### test helper

given the input

```gos
// file: input.gos
package example

import "testing"

func setUp() error { return nil }

func testSomething(tt testing.TB, input string) {
	require setUp()
}
```

the T2 output is

```
// golden-file: output.go
package example

import "testing"

func setUp() error { return nil }

func testSomething(tt testing.TB, input string) {
	if err := setUp(); err != nil {
	tt.Fatalf("%s: %v", "setUp()", err)
}
}
```

### bench

given the input

```gos
// file: input.gos
package example

import "testing"

func setUp() error { return nil }

func Benchmark(b *testing.B) {
	require setUp()
}
```

the T2 output is

```
// golden-file: output.go
package example

import "testing"

func setUp() error { return nil }

func Benchmark(b *testing.B) {
	if err := setUp(); err != nil {
	b.Fatalf("%s: %v", "setUp()", err)
}
}
```

### fuzz

given the input

```gos
// file: input.gos
package example

import "testing"

func setUp() error { return nil }

func Fuzz(f *testing.F) {
	require setUp()
}
```

the T2 output is

```
// golden-file: output.go
package example

import "testing"

func setUp() error { return nil }

func Fuzz(f *testing.F) {
	if err := setUp(); err != nil {
	f.Fatalf("%s: %v", "setUp()", err)
}
}
```