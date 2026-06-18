package main

import (
	"os"

	"nhatp.com/go/sugar/sdk"
	"nhatp.com/go/sugar/sugars/check"
)

func main() {
	if code := sdk.Run(os.Stdin, os.Stdout, os.Stderr, os.Args, check.New()); code != 0 {
		os.Exit(code)
	}
}
