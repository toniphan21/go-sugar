//go:build !dev

package main

import (
	"os"

	"nhatp.com/go/sugar/cli"
)

func main() {
	code := cli.Run(os.Stdin, os.Stdout, os.Stderr, os.Args)
	if code != 0 {
		os.Exit(code)
	}
}
