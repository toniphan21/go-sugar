## Semantic Transform - Message

given a go module

```go.mod
module github.com/you/repo

go 1.24
```

### prod doesn't support verb

given the input
```gos
// file: input.gos
package example

func setUp() error { return nil }

func DoSomething() {
	require setUp() "set up failed"
	require setUp() "%s failed"
	require setUp() "set up failed: %w"
	require setUp() "%s failed: %w"
	require setUp() "%%s failed: %%w"
}
```

the T2 output is

```
// golden-file: output.go
package example

func setUp() error { return nil }

func DoSomething() {
	if err := setUp() ; err != nil {
	panic("set up failed")
}
	if err := setUp() ; err != nil {
	panic("%s failed")
}
	if err := setUp() ; err != nil {
	panic("set up failed: %w")
}
	if err := setUp() ; err != nil {
	panic("%s failed: %w")
}
	if err := setUp() ; err != nil {
	panic("%%s failed: %%w")
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
	require setUp() "set up failed"
	require setUp() "%s failed"
	require setUp() "set up failed: %w"
	require setUp() "%s failed: %w"
	require setUp() "%%s failed: %%w"
}
```

the T2 output is

```
// golden-file: output.go
package example

import "testing"

func setUp() error { return nil }

func Test(t *testing.T) {
	if err := setUp() ; err != nil {
	t.Fatal("set up failed")
}
	if err := setUp() ; err != nil {
	t.Fatalf("%s failed", "setUp()")
}
	if err := setUp() ; err != nil {
	t.Fatalf("set up failed: %w", err)
}
	if err := setUp() ; err != nil {
	t.Fatalf("%s failed: %w", "setUp()", err)
}
	if err := setUp() ; err != nil {
	t.Fatal("%%s failed: %%w")
}
}
```
