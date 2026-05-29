package generatecmd

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"nhatp.com/go/gen-lib/cli"
	"nhatp.com/go/gen-lib/cli/color"
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/internal/util"
)

const cmdName = "generate"

type Arguments struct {
	WorkingDir string
	DryRun     bool
	Log        string
	LogLevel   slog.Level
}

func Run(stdin io.Reader, stdout io.Writer, stderr io.Writer, args Arguments) error {
	log, logFile, err := util.NewCLILogger(args.Log, stderr, args.LogLevel)
	if err != nil {
		return err
	}
	log = log.With("cmd", cmdName).WithGroup(cmdName)
	defer func() {
		log.Info(color.Generated("done"))
		if logFile != nil {
			if ce := logFile.Close(); ce != nil {
				err = ce
			}
		}
	}()

	if err = run(stdin, stdout, stderr, args, log); err != nil {
		log.Error(cli.ColorRed(err.Error()), slog.Any("error", err))
		return err
	}
	return nil
}

func run(stdin io.Reader, stdout, stderr io.Writer, args Arguments, log *slog.Logger) error {
	if args.DryRun {
		log.Info(color.Binary(sugar.BinaryName) + " " + color.Version(sugar.BinaryVersion) + " in DRY mode")
	} else {
		log.Info(color.Binary(sugar.BinaryName) + " " + color.Version(sugar.BinaryVersion))
	}

	workingDir, err := util.ResolveWorkingDir(args.WorkingDir)
	if err != nil {
		return err
	}
	log.Info(color.Binary(sugar.BinaryName) + " is working on directory: " + color.Input(workingDir))

	mod, err := sugar.NewModule(workingDir, sugar.DefaultConfig(), sugar.WithBinary(sugar.BinaryFullName), sugar.WithVersion(sugar.BinaryVersion))
	if err != nil {
		return err
	}

	files, err := mod.Generate()
	if err != nil {
		return err
	}

	for f, content := range files {
		if args.DryRun {
			_, _ = fmt.Fprintf(stdout, "// === go-sugar: %v    ===\n", f.SugarPath())
			_, _ = fmt.Fprintf(stdout, "// ---  desugar: %v ---\n", f.GoPath())
			_, _ = fmt.Fprint(stdout, string(content))
		} else {
			relPath := f.GoPath()
			absPath := filepath.Join(mod.Root, relPath)
			if err = os.WriteFile(absPath, content, 0644); err != nil {
				log.Error(cli.ColorRed(err.Error()))
				os.Exit(1)
			}
			log.Info(color.Binary(sugar.BinaryName) + " generated " + color.Generated(relPath) + " from " + color.Source(f.SugarPath()))
		}
	}
	return nil
}
