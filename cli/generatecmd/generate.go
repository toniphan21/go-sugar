package generatecmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"nhatp.com/go/gen-lib/cli"
	"nhatp.com/go/gen-lib/cli/color"
	"nhatp.com/go/sugar"
)

type Output struct {
	Argument   string       `json:"argument"`
	WorkingDir string       `json:"workingDir"`
	ModuleRoot string       `json:"moduleRoot"`
	Files      []OutputPair `json:"files"`
}

type OutputPair struct {
	Source    OutputFile `json:"source"`
	Generated OutputFile `json:"generated"`
}

type OutputFile struct {
	RelPath     string `json:"relPath"`
	DisplayPath string `json:"displayPath"`
	Content     string `json:"content"`
}

func runGenerate(stdin io.Reader, stdout, stderr io.Writer, args Arguments, log *slog.Logger) error {
	if args.DryRun {
		log.Info(color.Binary(sugar.BinaryName) + " " + color.Version(sugar.BinaryVersion) + " in DRY mode")
	} else {
		log.Info(color.Binary(sugar.BinaryName) + " " + color.Version(sugar.BinaryVersion))
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	mod, err := sugar.NewModule(wd, sugar.DefaultConfig(), sugar.WithBinary(sugar.BinaryFullName), sugar.WithVersion(sugar.BinaryVersion))
	if err != nil {
		return err
	}

	var outputs []Output

	targets, err := mod.Resolve(args.inputs()...)
	if err != nil {
		return err
	}

	for _, target := range targets {
		if target.IsDir {
			re := "recursively"
			if !target.Recursive {
				re = "(non-recursive)"
			}
			log.Info(color.Binary(sugar.BinaryName) + " is working on directory " + cli.ColorYellow(re) + ": " + color.Input(target.Path))
		}

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
			content, err := mod.GenerateFile(f)
			if err != nil {
				return err
			}

			source := OutputFile{
				DisplayPath: displayPathResolve(output.WorkingDir, output.ModuleRoot, f.SugarPath()),
				RelPath:     f.SugarPath(),
				Content:     string(f.Content()),
			}
			generated := OutputFile{
				DisplayPath: displayPathResolve(output.WorkingDir, output.ModuleRoot, f.GoPath()),
				RelPath:     f.GoPath(),
				Content:     string(content),
			}

			if args.DryRun || args.JSON {
				output.Files = append(output.Files, OutputPair{Source: source, Generated: generated})
				continue
			}

			absPath := filepath.Join(mod.Root, f.GoPath())
			if err = os.WriteFile(absPath, content, 0644); err != nil {
				log.Error(cli.ColorRed(err.Error()))
				os.Exit(1)
			}
			log.Info(color.Binary(sugar.BinaryName) + " generated " + color.Generated(generated.DisplayPath) + " from " + color.Source(source.DisplayPath))
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
				_, _ = fmt.Fprintf(stdout, "// === go-sugar: %v    ===\n", vv.Source.DisplayPath)
				_, _ = fmt.Fprintf(stdout, "// ---  desugar: %v ---\n", vv.Generated.DisplayPath)
				_, _ = fmt.Fprint(stdout, vv.Generated.Content)
			}
		}
	}
	return nil
}

func displayPathResolve(wd, root, path string) string {
	fp := filepath.Join(root, path)
	v, err := filepath.Rel(wd, fp)
	if err != nil {
		return path
	}
	return v
}
