package cli

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

const CodeExUsage = 64

type Runner[T any] func(stdin, stdout, stderr *os.File, args T) error

type UsageError interface {
	error

	UsageError()
}

func PrintHelp(stderr *os.File, help string, vars ...*bool) bool {
	for _, v := range vars {
		if *v {
			_, _ = fmt.Fprint(stderr, help)
			return true
		}
	}
	return false
}

func CmdUsage(out io.Writer, txt string) func() {
	return func() {
		_, _ = fmt.Fprint(out, txt)
	}
}

func InvokeRunner[T any](stdin, stdout, stderr *os.File, arg T, runner Runner[T], errorHandler func(error)) int {
	if err := runner(stdin, stdout, stderr, arg); err != nil {
		errorHandler(err)
		return 1
	}
	return 0
}

func PrintErrorTo(writer io.Writer) func(error) {
	return func(err error) {
		_, _ = fmt.Fprint(writer, err.Error())
	}
}

func PrintUsage(writer io.Writer, usage string) func(error) {
	return func(err error) {
		if _, ok := errors.AsType[UsageError](err); ok {
			_, _ = fmt.Fprint(writer, usage)
		}
	}
}

func IgnoreError(err error) {}

func LogLevel(verbosity *bool) slog.Level {
	v := slog.LevelInfo
	if *verbosity {
		v = slog.LevelDebug
	}
	return v
}

func FlagVal[T any](args ...*T) T {
	var v T
	for _, arg := range args {
		if arg != nil {
			v = *arg
		}
	}
	return v
}

func ReorderArgs(args []string, knownFlags ...string) []string {
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
