package fmtcmd

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"nhatp.com/go/gen-lib/cli"
	"nhatp.com/go/gen-lib/cli/color"
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/internal/util"
)

const cmdName = "fmt"

type Arguments struct {
	Args   []string
	DryRun bool
}

func (a *Arguments) inputs() []string {
	if len(a.Args) == 0 {
		return []string{"./"}
	}
	return a.Args
}

func Run(stdin io.Reader, stdout io.Writer, stderr io.Writer, args Arguments) error {
	log, logFile, err := util.NewCLILogger("", stderr, slog.LevelInfo)
	if logFile != nil {
		panic("log file is not allowed in fmt sub-command")
	}
	if err != nil {
		return err
	}
	defer log.Info(color.Generated("done"))

	if err = run(stdin, stdout, stderr, args, log); err != nil {
		log.Error(cli.ColorRed(err.Error()), slog.Any("error", err))
		return err
	}
	return nil
}

func run(stdin io.Reader, stdout, stderr io.Writer, args Arguments, log *slog.Logger) error {
	if args.DryRun {
		log.Info(color.Binary(sugar.BinaryName) + " " + color.Version(sugar.BinaryVersion) + " formatting in DRY mode")
	} else {
		log.Info(color.Binary(sugar.BinaryName) + " " + color.Version(sugar.BinaryVersion) + " formatting:")
	}

	wd, err := util.ResolveWorkingDir("")
	if err != nil {
		return err
	}

	mod, err := sugar.NewModule(wd, sugar.DefaultConfig())
	if err != nil {
		return err
	}

	fps, err := mod.Resolve(args.inputs()...)
	if err != nil {
		return err
	}

	for _, fp := range fps {
		content, err := mod.FormatFile(fp.RelPath)
		if err != nil {
			return err
		}

		if args.DryRun {
			_, _ = fmt.Fprintf(stdout, "// === go-sugar: %v ===\n", fp.DisplayPath)
			_, _ = fmt.Fprint(stdout, string(content))
		} else {
			if err = os.WriteFile(fp.AbsPath, content, 0644); err != nil {
				log.Error(cli.ColorRed(err.Error()))
				os.Exit(1)
			}
			log.Info("\t" + color.Source(fp.DisplayPath))
		}
	}
	return nil
}
