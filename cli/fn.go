package cli

import (
	"os"

	"nhatp.com/go/gen-lib/cli"
)

func disableColorIfNeeded(flag *bool, stdout *os.File) {
	isTerminal := false
	if stat, err := stdout.Stat(); err == nil {
		isTerminal = (stat.Mode() & os.ModeCharDevice) != 0
	}

	if *flag || os.Getenv("NO_COLOR") != "" || !isTerminal {
		cli.DisableColor()
	}
}
