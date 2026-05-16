package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"nhatp.com/go/gen-lib/cli"
	"nhatp.com/go/sugar/cli/lspcmd"
	"nhatp.com/go/sugar/cli/versioncmd"
)

const codeExUsage = 64

type Runner[T any] func(stdin io.Reader, stdout io.Writer, stderr io.Writer, args T) error

const usageText = `usage: go-sugar <command> [<args>...]

go-sugar - a superset of go with sugar.

See docs at https://nhatp.com/go/sugar

commands:
  generate   Generate Go code from gos file
  fmt        Format gos files
  lsp        Starts a language server for gos files
  info       Displays information about go-sugar environment
  version    Prints the version

`

func Run(stdin, stdout, stderr *os.File, args []string) int {
	if len(args) < 2 {
		_, _ = fmt.Fprint(stderr, usageText)
		return codeExUsage
	}

	switch args[1] {
	case "lsp":
		return lsp(stdin, stdout, stderr, args[2:], lspcmd.Run)
	case "version", "-version", "--version":
		return version(stdin, stdout, stderr, args[2:], versioncmd.Run)
	}

	return codeExUsage
}

// ---

const lspUsageText = `usage: go-sugar lsp [flags]

Prints the current version of go-sugar.

Flags:
  -log      The file to log LSP output to, or leave empty to disable logging.
  -h, -help Print this help message and exit.

`

func lsp(stdin, stdout, stderr *os.File, args []string, runner Runner[lspcmd.Arguments]) int {
	cmd := flag.NewFlagSet("lsp", flag.ContinueOnError)
	log := cmd.String("log", "", "The file to log LSP output to")
	h := cmd.Bool("h", false, "")
	help := cmd.Bool("help", false, "")

	cmd.Usage = func() {
		_, _ = fmt.Fprint(stderr, lspUsageText)
	}

	err := cmd.Parse(args)
	if err != nil {
		return codeExUsage
	}

	cli.DisableColor()
	if printHelp(stderr, lspUsageText, h, help) {
		return 0
	}

	arg := &lspcmd.Arguments{
		Log: *log,
	}

	if err = runner(stdin, stdout, stderr, *arg); err != nil {
		_, _ = fmt.Fprint(stderr, err.Error())
		return 1
	}
	return 0
}

// ---

const versionUsageText = `usage: go-sugar version [flags]

Prints the current version of go-sugar.

Flags:
  -v          Print the semantic version without the "v" prefix.
  -json       Print the version information in JSON format.
  -no-color   Disable color output.
  -h, -help   Print this help message and exit.

`

func version(stdin, stdout, stderr *os.File, args []string, runner Runner[versioncmd.Arguments]) int {
	cmd := flag.NewFlagSet("version", flag.ContinueOnError)
	noColor := cmd.Bool("no-color", false, "")
	h := cmd.Bool("h", false, "")
	help := cmd.Bool("help", false, "")
	semver := cmd.Bool("v", false, "")
	json := cmd.Bool("json", false, "")

	cmd.Usage = func() {
		_, _ = fmt.Fprint(stderr, versionUsageText)
	}

	err := cmd.Parse(args)
	if err != nil {
		return codeExUsage
	}

	disableColorIfNeeded(noColor, stdout)
	if printHelp(stderr, versionUsageText, h, help) {
		return 0
	}

	arg := &versioncmd.Arguments{}
	if *semver {
		arg.Format = versioncmd.FormatSemver
	}
	if *json {
		arg.Format = versioncmd.FormatJSON
	}

	if err = runner(stdin, stdout, stderr, *arg); err != nil {
		_, _ = fmt.Fprint(stderr, err.Error())
		return 1
	}
	return 0
}
