package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"nhatp.com/go/gen-lib/cli"
	"nhatp.com/go/sugar/cli/generatecmd"
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
	case "generate":
		return generate(stdin, stdout, stderr, args[2:], generatecmd.Run)
	case "lsp":
		return lsp(stdin, stdout, stderr, args[2:], lspcmd.Run)
	case "version", "-version", "--version":
		return version(stdin, stdout, stderr, args[2:], versioncmd.Run)
	}

	return codeExUsage
}

// ---

const generateUsageText = `usage: go-sugar generate [flags]

Generates Go code from go-sugar files.

Flags:
  -log              The file to log the command output to, or leave empty to disable logging.
  -w, -working-dir  The working directory to use when generating the code. (default ".")
  -d, -dry          Preview changes without writing to disk.
  -v                Set log verbosity level to "debug". (default "info")
  -no-color         Disable color output.
  -h, -help         Print this help message and exit.

`

func generate(stdin, stdout, stderr *os.File, args []string, runner Runner[generatecmd.Arguments]) int {
	cmd := flag.NewFlagSet("lsp", flag.ContinueOnError)
	log := cmd.String("log", "", "The file to log the command output to")

	w := cmd.String("w", ".", "The working directory to use when generating the code")
	workingDir := cmd.String("working-dir", ".", "The working directory to use when generating the code")

	d := cmd.Bool("d", false, "Preview changes without writing to disk")
	dry := cmd.Bool("dry", false, "Preview changes without writing to disk")

	verbosity := cmd.Bool("v", false, "Set log verbosity level to debug")
	noColor := cmd.Bool("no-color", false, "")
	h := cmd.Bool("h", false, "")
	help := cmd.Bool("help", false, "")

	cmd.Usage = cmdUsage(stderr, generateUsageText)
	if err := cmd.Parse(args); err != nil {
		return codeExUsage
	}

	disableColorIfNeeded(noColor, stdout)
	if printHelp(stderr, generateUsageText, h, help) {
		return 0
	}

	arg := generatecmd.Arguments{
		WorkingDir: flagVal(w, workingDir),
		DryRun:     *d || *dry,
		Log:        flagVal(log),
		LogLevel:   logLevel(verbosity),
	}
	return invokeRunner(stdin, stdout, stderr, arg, runner)
}

// ---

const lspUsageText = `usage: go-sugar lsp [flags]

Starts a language server for go-sugar files.

Flags:
  -log       The file to log LSP output to, or leave empty to disable logging.
  -v         Set log verbosity level to "debug". (default "info")
  -h, -help  Print this help message and exit.

`

func lsp(stdin, stdout, stderr *os.File, args []string, runner Runner[lspcmd.Arguments]) int {
	cmd := flag.NewFlagSet("lsp", flag.ContinueOnError)
	log := cmd.String("log", "", "The file to log LSP output to")
	verbosity := cmd.Bool("v", false, "Set log verbosity level to debug")
	h := cmd.Bool("h", false, "")
	help := cmd.Bool("help", false, "")

	cmd.Usage = cmdUsage(stderr, lspUsageText)
	if err := cmd.Parse(args); err != nil {
		return codeExUsage
	}

	cli.DisableColor()
	if printHelp(stderr, lspUsageText, h, help) {
		return 0
	}

	arg := lspcmd.Arguments{
		Log:      *log,
		LogLevel: logLevel(verbosity),
	}
	return invokeRunner(stdin, stdout, stderr, arg, runner)
}

// ---

const versionUsageText = `usage: go-sugar version [flags]

Prints the current version of go-sugar.

Flags:
  -v         Print the semantic version without the "v" prefix.
  -json      Print the version information in JSON format.
  -no-color  Disable color output.
  -h, -help  Print this help message and exit.

`

func version(stdin, stdout, stderr *os.File, args []string, runner Runner[versioncmd.Arguments]) int {
	cmd := flag.NewFlagSet("version", flag.ContinueOnError)
	semver := cmd.Bool("v", false, "")
	json := cmd.Bool("json", false, "")
	noColor := cmd.Bool("no-color", false, "")
	h := cmd.Bool("h", false, "")
	help := cmd.Bool("help", false, "")

	cmd.Usage = cmdUsage(stderr, versionUsageText)
	if err := cmd.Parse(args); err != nil {
		return codeExUsage
	}

	disableColorIfNeeded(noColor, stdout)
	if printHelp(stderr, versionUsageText, h, help) {
		return 0
	}

	arg := versioncmd.Arguments{}
	if *semver {
		arg.Format = versioncmd.FormatSemver
	}
	if *json {
		arg.Format = versioncmd.FormatJSON
	}
	return invokeRunner(stdin, stdout, stderr, arg, runner)
}
