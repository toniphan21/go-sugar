package cli

import (
	"fmt"
	"io"
	"log/slog"
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

func printHelp(stderr *os.File, help string, vars ...*bool) bool {
	for _, v := range vars {
		if *v {
			_, _ = fmt.Fprint(stderr, help)
			return true
		}
	}
	return false
}

func cmdUsage(out io.Writer, txt string) func() {
	return func() {
		_, _ = fmt.Fprint(out, txt)
	}
}

func invokeRunner[T any](stdin, stdout, stderr *os.File, arg T, runner Runner[T]) int {
	if err := runner(stdin, stdout, stderr, arg); err != nil {
		_, _ = fmt.Fprint(stderr, err.Error())
		return 1
	}
	return 0
}

func logLevel(verbosity *bool) slog.Level {
	v := slog.LevelInfo
	if *verbosity {
		v = slog.LevelDebug
	}
	return v
}

func flagVal[T any](args ...*T) T {
	var v T
	for _, arg := range args {
		if arg != nil {
			v = *arg
		}
	}
	return v
}
