//go:build dev

package main

import (
	"os"

	"nhatp.com/go/sugar/cli"
)

func main() {
	args := os.Args
	args = []string{"go-sugar", "test", "-s", "./sugars/check/testdata/generate/basic.md"}

	code := cli.Run(os.Stdin, os.Stdout, os.Stderr, args)
	if code != 0 {
		os.Exit(code)
	}
}
