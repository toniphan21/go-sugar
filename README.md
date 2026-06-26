> **Status: early development**. Core architecture is settled and the first reference plugins (require, assert) are
> landing. This is a work-in-progress published openly to share the design.

## Introduction

`go-sugar` is a source-to-source transformation toolchain that extends Go with syntactic sugar. You write `.gos` files -
a superset of Go with sugar constructs - and `go-sugar` compiles them down to plain, idiomatic `.go` that the standard
toolchain compiles, formats, and lints without modification.

The goal is to **cut boilerplate without leaving the Go ecosystem**. The near-term focus is test code, where boilerplate
is densest. The longer-term vision includes a plugin protocol, LSP integration, and eventually defensive-programming
constructs usable in production (think C's `assert.h` with `NDEBUG`).

## Example

The first two reference sugar plugins, `require` and `assert`, target the most repetitive parts of Go tests.

### `require`: keyword for error-handling boilerplate

test file with `require` sugar:

```gos
// test.gos
package demo

import "testing"

func Test_Something(t *testing.T) {
	require setUp()
	svc := require NewService() "cannot make Service"
	// ...
}
```

will generate `test.gos.go`:

```go
package demo

import "testing"

func Test_Something(t *testing.T) {
	if err := setUp(); err != nil {
		t.Fatalf("%s: %v", "setUp()", err)
	}
	svc, err := NewService()
	if err != nil {
		t.Fatal("cannot make Service")
	}
	// ...
}
```

### `assert`: built-in assert function with informative failure messages

test file with `assert` sugar:

```go
// test.gos

package demo

import "testing"

func Test_Something(t *testing.T) { 
	// ... 
	
	assert(age > 18)
	assert(username == "test")
	assert(result == expected)
	assert(a == "a" || !b)
	assert(called != 0, "service is not called")
	assert(returned == 0, "program exits with error code %v", returned)
}
```

will generate `test.gos.go`:

```go
package demo

import "testing"

func Test_Something(t *testing.T) {
	// ...

	if !(age > 18) {
		t.Errorf("%s\n\t got: %#v", "age > 18", age)
	}
	if !(username == "test") {
		t.Errorf("%s\n\t got: %#v", "username == \"test\"", username)
	}
	if !(result == expected) {
		t.Errorf("%s\n\t got: %#v\n\twant: %#v", "result == expected", result, expected)
	}
	if !(a == "a" || !b) {
		t.Error("assertion failed: a == \"a\" || !b")
	}
	if !(called != 0) {
		t.Error("service is not called")
	}
	if !(returned == 0) {
		t.Errorf("program exits with error code %v", returned)
	}
}
```

## Status

- [x] Infrastructure:
    - [x] Lexical parser helper
    - [x] Golden tests (use .md as virtual file system)
    - [x] Sugar Plugin API
    - [x] SDK and transport to independent sugar binary
- [ ] Backend - *in progress*
    - [x] Generate pipeline with watch
    - [x] Format pipeline with watch
    - [ ] Built-in sugars - *in progress*
        - [x] require
        - [ ] assert - *in progress*
    - [ ] Test pipeline with watch and sourcemap rewrite - *in progress*
    - [ ] LSP proxying `gopls` with sourcemap rewrite
- [ ] Frontend:
    - [ ] LSP and example configs for nvim
    - [ ] vscode extension
    - [ ] IDEA plugin

## Tools

- [lexeme-viewer](tools/lexeme-viewer): wasm tool to view lexeme in browser, use port 39800
- [railroad-diagram](tools/railroad-diagram): tool to generate railroad diagram svg, use port 39801
  (thanks to [tabatkins/railroad-diagrams](https://github.com/tabatkins/railroad-diagrams/))

## Contributing & License

PRs are welcome! Distributed under the Apache License 2.0.

