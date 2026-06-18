package main

import (
	"os"

	"nhatp.com/go/sugar/sdk"
	"nhatp.com/go/sugar/sugars/require"
)

func main() {
	if code := sdk.Run(os.Stdin, os.Stdout, os.Stderr, os.Args, require.New()); code != 0 {
		os.Exit(code)
	}
}
