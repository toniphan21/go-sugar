//go:build dev

package main

import (
	"os"

	"nhatp.com/go/sugar/cli"
)

func main() {
	args := os.Args
	args = []string{"go-sugar", "gen", "-w", "./example/hellos.gos"}

	code := cli.Run(os.Stdin, os.Stdout, os.Stderr, args)
	if code != 0 {
		os.Exit(code)
	}
}
