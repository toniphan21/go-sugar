package fmtcmd

import (
	"encoding/json"
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
	JSON   bool
}

func (a *Arguments) inputs() []string {
	if len(a.Args) == 0 {
		return []string{"./"}
	}
	return a.Args
}

func Run(stdin, stdout, stderr *os.File, args Arguments) error {
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

type Output struct {
	Argument   string       `json:"argument"`
	WorkingDir string       `json:"workingDir"`
	ModuleRoot string       `json:"moduleRoot"`
	Files      []OutputFile `json:"files"`
}

type OutputFile struct {
	RelPath     string `json:"relPath"`
	DisplayPath string `json:"displayPath"`
	Original    string `json:"original"`
	Formatted   string `json:"formatted"`
}

func run(stdin io.Reader, stdout, stderr io.Writer, args Arguments, log *slog.Logger) error {
	if args.DryRun {
		log.Info(color.Binary(sugar.BinaryName) + " " + color.Version(sugar.BinaryVersion) + " formatting in DRY mode")
	} else {
		log.Info(color.Binary(sugar.BinaryName) + " " + color.Version(sugar.BinaryVersion) + " formatting:")
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	mod, err := sugar.NewModule(wd, sugar.DefaultConfig())
	if err != nil {
		return err
	}

	var outputs []Output

	targets, err := mod.Resolve(args.inputs()...)
	if err != nil {
		return err
	}

	for _, target := range targets {
		output := Output{
			Argument:   target.Input,
			WorkingDir: target.WorkingDir,
			ModuleRoot: target.Root,
		}

		fps, err := target.Resolve()
		if err != nil {
			return err
		}

		for _, fp := range fps {
			f, ok := mod.File(fp.RelPath)
			if !ok {
				return fmt.Errorf("file not found: %v", fp.RelPath)
			}

			content, err := mod.FormatFile(f)
			if err != nil {
				return err
			}

			outputFile := OutputFile{
				DisplayPath: fp.DisplayPath,
				RelPath:     fp.RelPath,
				Original:    string(f.Content()),
				Formatted:   string(content),
			}

			if args.DryRun || args.JSON {
				output.Files = append(output.Files, outputFile)
				continue
			}

			if err = os.WriteFile(fp.AbsPath, content, 0644); err != nil {
				log.Error(cli.ColorRed(err.Error()))
				os.Exit(1)
			}
			log.Info("\t" + color.Source(fp.DisplayPath))
		}
		outputs = append(outputs, output)
	}

	switch {
	case args.JSON:
		out, err := json.Marshal(outputs)
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintln(stdout, string(out))
	case args.DryRun:
		for _, v := range outputs {
			for _, vv := range v.Files {
				_, _ = fmt.Fprintf(stdout, "// === go-sugar: %v    ===\n", vv.DisplayPath)
				_, _ = fmt.Fprint(stdout, vv.Formatted)
			}
		}
	}
	return nil
}
