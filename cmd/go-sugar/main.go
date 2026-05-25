//go:build !dev

package main

import (
	"os"
)

func main() {
	if code := cli.Run(os.Stdin, os.Stdout, os.Stderr, os.Args); code != 0 {
		os.Exit(code)
	}
}
