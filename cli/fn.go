package cli

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"nhatp.com/go/gen-lib/cli"
)

type UsageError interface {
	error

	UsageError()
}

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

func invokeRunner[T any](stdin, stdout, stderr *os.File, arg T, runner Runner[T], errorHandler func(error)) int {
	if err := runner(stdin, stdout, stderr, arg); err != nil {
		errorHandler(err)
		return 1
	}
	return 0
}

func printErrorTo(writer io.Writer) func(error) {
	return func(err error) {
		_, _ = fmt.Fprint(writer, err.Error())
	}
}

func printUsage(writer io.Writer, usage string) func(error) {
	return func(err error) {
		if _, ok := errors.AsType[UsageError](err); ok {
			_, _ = fmt.Fprint(writer, usage)
		}
	}
}

func ignoreError(err error) {}

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

func reorderArgs(args []string, knownFlags ...string) []string {
	known := make(map[string]bool)
	for _, v := range knownFlags {
		known["-"+v] = true
		known["--"+v] = true
	}

	isKnownFlag := func(v string) bool {
		if known[v] {
			return true
		}
		for k := range known {
			if strings.HasPrefix(v, k+"=") {
				return true
			}
		}
		return false
	}

	var flags []string
	var positional []string
	for _, v := range args {
		if isKnownFlag(v) {
			flags = append(flags, v)
		} else {
			positional = append(positional, v)
		}
	}
	return append(flags, positional...)
}
