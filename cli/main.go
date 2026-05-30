package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"nhatp.com/go/gen-lib/cli"
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/cli/fmtcmd"
	"nhatp.com/go/sugar/cli/generatecmd"
	"nhatp.com/go/sugar/cli/goldencmd"
	"nhatp.com/go/sugar/cli/lspcmd"
	"nhatp.com/go/sugar/cli/versioncmd"
	"nhatp.com/go/sugar/sugars/check"
)

const codeExUsage = 64

type Runner[T any] func(stdin io.Reader, stdout io.Writer, stderr io.Writer, args T) error

const usageText = `usage: go-sugar <command> [<args>...]

go-sugar - a superset of go with sugar.

See docs at https://nhatp.com/go/sugar

commands:
  fmt       Format gos files
  lsp       Starts a language server for gos files
  info      Displays information about go-sugar environment
  generate  Generate Go code from gos file
  golden    Run Markdown golden tests
  version   Prints the version

`

func Run(stdin, stdout, stderr *os.File, args []string) int {
	if len(args) < 2 {
		_, _ = fmt.Fprint(stderr, usageText)
		return codeExUsage
	}

	registerPlugins()
	switch args[1] {
	case "fmt", "format":
		return format(stdin, stdout, stderr, args[2:], fmtcmd.Run)
	case "gen", "generate":
		return generate(stdin, stdout, stderr, args[2:], generatecmd.Run)
	case "lsp":
		return lsp(stdin, stdout, stderr, args[2:], lspcmd.Run)
	case "golden":
		return golden(stdin, stdout, stderr, args[2:], goldencmd.Run)
	case "version", "-version", "--version":
		return version(stdin, stdout, stderr, args[2:], versioncmd.Run)
	}

	return codeExUsage
}

func registerPlugins() {
	sugar.Register(check.New())
}

// ---

const fmtUsageText = `usage: go-sugar fmt [flags] [FILE|DIR|PATTERN...]

Format go-sugar files.

  ./     Format go-sugar files in the current directory.
  ./...  Format go-sugar files in the current directory tree, recursively.

If no arguments are given, defaults to the current directory (non-recursive).

Flags:
  -d, -dry   Preview changes without writing to disk.
  -no-color  Disable color output.
  -h, -help  Print this help message and exit.

`

func format(stdin, stdout, stderr *os.File, args []string, runner Runner[fmtcmd.Arguments]) int {
	cmd := flag.NewFlagSet("lsp", flag.ContinueOnError)
	d := cmd.Bool("d", false, "Preview changes without writing to disk")
	dry := cmd.Bool("dry", false, "Preview changes without writing to disk")

	noColor := cmd.Bool("no-color", false, "")
	h := cmd.Bool("h", false, "")
	help := cmd.Bool("help", false, "")

	cmd.Usage = cmdUsage(stderr, fmtUsageText)
	reorderableFlags := []string{
		"d", "dry", "no-color",
	}
	if err := cmd.Parse(reorderArgs(args, reorderableFlags...)); err != nil {
		return codeExUsage
	}

	disableColorIfNeeded(noColor, stdout)
	if printHelp(stderr, generateUsageText, h, help) {
		return 0
	}

	arg := fmtcmd.Arguments{
		Args:   cmd.Args(),
		DryRun: *d || *dry,
	}
	return invokeRunner(stdin, stdout, stderr, arg, runner, printUsage(stderr, generateUsageText))
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
	reorderableFlags := []string{
		"w", "working-dir", "d", "dry", "v", "no-color",
	}
	if err := cmd.Parse(reorderArgs(args, reorderableFlags...)); err != nil {
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
	return invokeRunner(stdin, stdout, stderr, arg, runner, printUsage(stderr, generateUsageText))
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

const goldenUsageText = `usage: go-sugar golden [flags] FILE [FILE ...]

Run Markdown golden tests. If no pipeline flag is given, defaults to -gen.

Arguments:
  FILE              Markdown test files to run.

Flags:
  -t1,  -structural  Run T1 StructuralTransform test.
  -t2,  -semantic    Run T2 SemanticTransform test.
  -t3,  -restore     Run T3 RestoreTransform test.
  -fmt, -format      Run format pipeline test (T1 + gofmt + T3).
  -gen, -generate    Run generate pipeline test (T1 + T2 + gofmt).
  -n, -name          Filter tests by name. (case insensitive)
  -s, -show-setup    Show test setup steps. (default: false)
  -t, -tab-size      Number of spaces per tab. (default: 8)
  -e, -emit-code     Emit code if the test passes. If empty, looks for path in a Markdown comment.
  -log               The file to log the command output to, or leave empty to disable logging.
  -v                 Set log verbosity level to "debug". (default "info")
  -no-color          Disable color output.
  -h, -help          Print this help message and exit.

`

func golden(stdin, stdout, stderr *os.File, args []string, runner Runner[goldencmd.Arguments]) int {
	cmd := flag.NewFlagSet("lsp", flag.ContinueOnError)
	t1s := cmd.Bool("t1", false, "")
	t1l := cmd.Bool("structural", false, "")

	t2s := cmd.Bool("t2", false, "")
	t2l := cmd.Bool("semantic", false, "")

	t3s := cmd.Bool("t3", false, "")
	t3l := cmd.Bool("restore", false, "")

	tfs := cmd.Bool("fmt", false, "")
	tfl := cmd.Bool("format", false, "")

	tgs := cmd.Bool("gen", false, "")
	tgl := cmd.Bool("generate", false, "")

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

	cmd.Usage = cmdUsage(stderr, goldenUsageText)
	reorderableFlags := []string{
		"t1", "structural", "t2", "semantic", "t3", "restore", "fmt", "format", "gen", "generate",
		"n", "name", "s", "show-setup", "t", "tab-size", "e", "emit-code", "log", "v", "no-color",
	}
	if err := cmd.Parse(reorderArgs(args, reorderableFlags...)); err != nil {
		return codeExUsage
	}

	disableColorIfNeeded(noColor, stdout)
	if printHelp(stderr, goldenUsageText, h, help) {
		return 0
	}

	arg := goldencmd.Arguments{
		Files:     cmd.Args(),
		Name:      flagVal(n, name),
		ShowSetup: *s || *showSetup,
		TabSize:   flagVal(t, tabSize),
		EmitCode:  flagVal(e, emitCode),
		Log:       *log,
		LogLevel:  logLevel(verbosity),
	}
	switch {
	case *t1s || *t1l:
		arg.Type = goldencmd.TypeStructuralTransform
	case *t2s || *t2l:
		arg.Type = goldencmd.TypeSemanticTransform
	case *t3s || *t3l:
		arg.Type = goldencmd.TypeRestoreTransform
	case *tfs || *tfl:
		arg.Type = goldencmd.TypeFormatPipeline
	case *tgs || *tgl:
		arg.Type = goldencmd.TypeGeneratePipeline
	default:
		arg.Type = goldencmd.TypeGeneratePipeline
	}
	return invokeRunner(stdin, stdout, stderr, arg, runner, printUsage(stderr, goldenUsageText))
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
