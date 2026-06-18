package sdk

import (
	"flag"
	"fmt"
	"os"

	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/internal/cli"
	"nhatp.com/go/sugar/sdk/startcmd"
	"nhatp.com/go/sugar/sdk/versioncmd"
)

func usageText(sugar sugar.Sugar) string {
	t := `usage: %v <command> [<args>...]

%v - %v

commands:
  start    Start plugin server
  version  Print the version

`

	binary := sugar.Binary()
	return fmt.Sprintf(t, binary.Name, binary.Name, binary.Usage)
}

func Run(stdin, stdout, stderr *os.File, args []string, sugar sugar.Sugar) int {
	if len(args) < 2 {
		_, _ = fmt.Fprint(stderr, usageText(sugar))
		return cli.CodeExUsage
	}

	switch args[1] {
	case "start":
		return start(stdin, stdout, stderr, args[2:], sugar, startcmd.Run)
	case "version", "-version", "--version":
		return version(stdin, stdout, stderr, args[2:], sugar, versioncmd.Run)
	}

	return cli.CodeExUsage
}

// ---

func startUsageText(sugar sugar.Sugar) string {
	t := `usage: %v start [flags]

Starts a server for go-sugar plugin.

Flags:
  -log FILE  The file to log the command output to.
  -v         Set log verbosity level to "debug". (default "info")
  -h, -help  Print this help message and exit.

`
	return fmt.Sprintf(t, sugar.Binary().Name)
}

func start(stdin, stdout, stderr *os.File, args []string, sugar sugar.Sugar, runner cli.Runner[startcmd.Arguments]) int {
	cmd := flag.NewFlagSet("lsp", flag.ContinueOnError)
	log := cmd.String("log", "", "The file to log output to")
	verbosity := cmd.Bool("v", false, "Set log verbosity level to debug")
	h := cmd.Bool("h", false, "")
	help := cmd.Bool("help", false, "")

	cmd.Usage = cli.CmdUsage(stderr, startUsageText(sugar))
	if err := cmd.Parse(args); err != nil {
		return cli.CodeExUsage
	}

	if cli.PrintHelp(stderr, startUsageText(sugar), h, help) {
		return 0
	}

	arg := startcmd.Arguments{
		Sugar:    sugar,
		Log:      *log,
		LogLevel: cli.LogLevel(verbosity),
	}
	return cli.InvokeRunner(stdin, stdout, stderr, arg, runner, cli.PrintErrorTo(stderr))
}

// ---

func versionUsageText(sugar sugar.Sugar) string {
	t := `usage: %v version [flags]

Prints the current version of %v.

Flags:
  -v         Print the semantic version without the "v" prefix.
  -json      Print the version information in JSON format.
  -h, -help  Print this help message and exit.

`
	binary := sugar.Binary()
	return fmt.Sprintf(t, binary.Name, binary.Name)
}

func version(stdin, stdout, stderr *os.File, args []string, sugar sugar.Sugar, runner cli.Runner[versioncmd.Arguments]) int {
	cmd := flag.NewFlagSet("version", flag.ContinueOnError)
	semver := cmd.Bool("v", false, "")
	json := cmd.Bool("json", false, "")
	h := cmd.Bool("h", false, "")
	help := cmd.Bool("help", false, "")

	cmd.Usage = cli.CmdUsage(stderr, versionUsageText(sugar))
	if err := cmd.Parse(args); err != nil {
		return cli.CodeExUsage
	}

	if cli.PrintHelp(stderr, versionUsageText(sugar), h, help) {
		return 0
	}

	arg := versioncmd.Arguments{Sugar: sugar}
	if *semver {
		arg.Format = versioncmd.FormatSemver
	}
	if *json {
		arg.Format = versioncmd.FormatJSON
	}
	return cli.InvokeRunner(stdin, stdout, stderr, arg, runner, cli.PrintErrorTo(stderr))
}
