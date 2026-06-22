## Semantic Transform - Identifiers

given a go module

```go.mod
module github.com/you/repo

go 1.24
```

### no identifiers

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

### with identifiers

given the input

```gos
// file: input.gos
package example

import "testing"

func setUp() error { return nil }

func Test(t *testing.T) {
	x := require setUp()
}
```

the T2 output is

```
// golden-file: output.go
package example

import "testing"

func setUp() error { return nil }

func Test(t *testing.T) {
	x, err := setUp()
if err != nil {
	t.Fatalf("%s: %v", "setUp()", err)
}
}
```