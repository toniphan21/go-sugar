package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"nhatp.com/go/gen-lib/cli"
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/cli/generatecmd"
	"nhatp.com/go/sugar/cli/lspcmd"
	"nhatp.com/go/sugar/cli/testcmd"
	"nhatp.com/go/sugar/cli/versioncmd"
	"nhatp.com/go/sugar/sugars/check"
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

	registerPlugins()
	switch args[1] {
	case "generate":
		return generate(stdin, stdout, stderr, args[2:], generatecmd.Run)
	case "lsp":
		return lsp(stdin, stdout, stderr, args[2:], lspcmd.Run)
	case "test":
		return test(stdin, stdout, stderr, args[2:], testcmd.Run)
	case "version", "-version", "--version":
		return version(stdin, stdout, stderr, args[2:], versioncmd.Run)
	}

	return codeExUsage
}

func registerPlugins() {
	sugar.Register(check.New())
}

// ---

const generateUsageText = `usage: go-sugar generate [flags]

Generates Go code from go-sugar files.

Flags:
  -w, -working-dir  The working directory to use when generating the code. (default ".")
  -d, -dry          Preview changes without writing to disk.
  -log              The file to log the command output to, or leave empty to disable logging.
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
	return invokeRunner(stdin, stdout, stderr, arg, runner, ignoreError)
}

// ---

const lspUsageText = `usage: go-sugar lsp [flags]

Starts a language server for go-sugar files.

Flags:
  -log       The file to log the command output to, or leave empty to disable logging.
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
	return invokeRunner(stdin, stdout, stderr, arg, runner, printErrorTo(stderr))
}

// ---

const testUsageText = `usage: go-sugar test [flags] FILE [FILE ...]

Run Markdown golden tests. By default, tests the "generate" pipeline.

Arguments:
  FILE              Markdown test files to run.

Flags:
  -t1, -structural  Run T1 StructuralTransform test.
  -t2, -semantic    Run T2 SemanticTransform test.
  -t3, -restore     Run T3 RestoreTransform test.
  -n, -name         Run tests matching a name. (case insensitive)
  -s, -show-setup   Show test setup steps. (default: false)
  -t, -tab-size     Number of spaces per tab. (default: 8)
  -e, -emit-code    Emit code if the test passes. If empty, looks for path in a Markdown comment.
  -log              The file to log the command output to, or leave empty to disable logging.
  -v                Set log verbosity level to "debug". (default "info")
  -no-color         Disable color output.
  -h, -help         Print this help message and exit.

`

func test(stdin, stdout, stderr *os.File, args []string, runner Runner[testcmd.Arguments]) int {
	cmd := flag.NewFlagSet("lsp", flag.ContinueOnError)
	t1 := cmd.Bool("t1", false, "")
	t2 := cmd.Bool("t2", false, "")
	t3 := cmd.Bool("t3", false, "")

	n := cmd.String("n", "", "")
	name := cmd.String("name", "", "")

	s := cmd.Bool("s", false, "")
	showSetup := cmd.Bool("show-setup", false, "")

	t := cmd.Int("t", 8, "")
	tabSize := cmd.Int("tab-size", 8, "")

	e := cmd.String("e", "", "")
	emitCode := cmd.String("emit-code", "", "")

	log := cmd.String("log", "", "")
	verbosity := cmd.Bool("v", false, "")
	noColor := cmd.Bool("no-color", false, "")
	h := cmd.Bool("h", false, "")
	help := cmd.Bool("help", false, "")

	cmd.Usage = cmdUsage(stderr, testUsageText)
	if err := cmd.Parse(args); err != nil {
		return codeExUsage
	}

	disableColorIfNeeded(noColor, stdout)
	if printHelp(stderr, testUsageText, h, help) {
		return 0
	}

	arg := testcmd.Arguments{
		Files:     cmd.Args(),
		Name:      flagVal(n, name),
		ShowSetup: *s || *showSetup,
		TabSize:   flagVal(t, tabSize),
		EmitCode:  flagVal(e, emitCode),
		Log:       *log,
		LogLevel:  logLevel(verbosity),
	}
	switch {
	case *t1:
		arg.Type = testcmd.TypeStructuralTransform
	case *t2:
		arg.Type = testcmd.TypeSemanticTransform
	case *t3:
		arg.Type = testcmd.TypeRestoreTransform
	default:
		arg.Type = testcmd.TypeGenerate
	}
	return invokeRunner(stdin, stdout, stderr, arg, runner, ignoreError)
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
	return invokeRunner(stdin, stdout, stderr, arg, runner, printErrorTo(stderr))
}
