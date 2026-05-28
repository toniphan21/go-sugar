package goldencmd

import (
	"fmt"
	"io"
	"log/slog"
	"path/filepath"

	genlib "nhatp.com/go/gen-lib"
	"nhatp.com/go/gen-lib/cli"
	"nhatp.com/go/gen-lib/cli/color"
	"nhatp.com/go/gen-lib/file"
	"nhatp.com/go/gen-lib/gentest"
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/internal/util"
	"nhatp.com/go/sugar/sugartest"
)

const cmdName = "test"

type errNoFileSpecified struct {
}

func (e *errNoFileSpecified) UsageError() {}

func (e *errNoFileSpecified) Error() string {
	return "no files specified"
}

type Type int

const (
	TypeGenerate Type = iota
	TypeStructuralTransform
	TypeSemanticTransform
	TypeRestoreTransform
)

type Arguments struct {
	Files     []string
	Type      Type
	Name      string
	ShowSetup bool
	TabSize   int
	EmitCode  string
	Log       string
	LogLevel  slog.Level
}

func Run(stdin io.Reader, stdout io.Writer, stderr io.Writer, args Arguments) error {
	log, logFile, err := util.NewCLILogger(args.Log, stderr, args.LogLevel)
	if err != nil {
		return err
	}
	log = log.With("cmd", cmdName).WithGroup(cmdName)
	defer func() {
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
	if len(args.Files) == 0 {
		return &errNoFileSpecified{}
	}

	cmd := &cli.TestRunner{
		Logger: log,
		FilePathResolver: cli.WithVanityURLFilePathResolver(map[string]string{
			"repo://":      sugar.RawGitRefsURL + "/heads/main/",
			"repo-refs://": sugar.RawGitRefsURL,
		}),
		Files:     args.Files,
		Name:      args.Name,
		TabSize:   args.TabSize,
		ShowSetup: args.ShowSetup,
		EmitPath:  args.EmitCode,
	}

	switch args.Type {
	case TypeGenerate:
		runGenerate(cmd)
	case TypeStructuralTransform:
		runStructuralTransform(cmd)
	case TypeSemanticTransform:
		runSemanticTransform(cmd)
	case TypeRestoreTransform:
		runRestoreTransform(cmd)
	}
	return nil
}

func runGenerate(cmd *cli.TestRunner) {
	cmd.Print("Running tests with " + color.Binary(sugar.BinaryName) + " " + color.Version(sugar.BinaryVersion))
	cmd.Print("")
	cmd.RunTestCase = func(tc cli.TestCase, options map[string]any) (genlib.FileManager, error) {
		dir := tc.TestDir
		if err := genlib.SetupSourceCode(dir, tc.SourceFiles); err != nil {
			return nil, err
		}

		mod, err := sugar.NewModule(dir, sugar.Config{})
		if err != nil {
			return nil, err
		}

		files, err := mod.Generate()
		if err != nil {
			return nil, err
		}

		fm := genlib.NewFileManager(dir)
		for f, content := range files {
			relPath := f.GoPath()
			absPath := filepath.Join(dir, relPath)
			if err = fm.Add(&file.GoFile{Path: absPath, Content: string(content)}); err != nil {
				return nil, err
			}
		}
		return fm, nil
	}
	cmd.Run()
}

func runStructuralTransform(cmd *cli.TestRunner) {
	cmd.Print("Running " + color.Source("T1 - StructuralTransform") + " tests with " + color.Binary(sugar.BinaryName) + " " + color.Version(sugar.BinaryVersion))
	cmd.Print("")
	cmd.RunTestCase = func(tc cli.TestCase, options map[string]any) (genlib.FileManager, error) {
		return runTransformTest(tc, "output.go", sugartest.PerformStructuralTransform)
	}
	cmd.Run()
}

func runSemanticTransform(cmd *cli.TestRunner) {
	cmd.Print("Running " + color.Source("T2 - SemanticTransform") + " tests with " + color.Binary(sugar.BinaryName) + " " + color.Version(sugar.BinaryVersion))
	cmd.Print("")
	cmd.RunTestCase = func(tc cli.TestCase, options map[string]any) (genlib.FileManager, error) {
		return runTransformTest(tc, "output.go", sugartest.PerformSemanticTransform)
	}
	cmd.Run()
}

func runRestoreTransform(cmd *cli.TestRunner) {
	cmd.Print("Running " + color.Source("T3 - RestoreTransform") + " tests with " + color.Binary(sugar.BinaryName) + " " + color.Version(sugar.BinaryVersion))
	cmd.Print("")
	cmd.RunTestCase = func(tc cli.TestCase, options map[string]any) (genlib.FileManager, error) {
		return runTransformTest(tc, "output.gos", sugartest.PerformRestoreTransform)
	}
	cmd.Run()
}

func runTransformTest(tc cli.TestCase, outFileName string, transform func(string, sugar.Config, gentest.MarkdownTestCase) ([]byte, error)) (genlib.FileManager, error) {
	dir := tc.TestDir
	mtc := gentest.MarkdownTestCase{
		SourceFiles: tc.SourceFiles,
	}

	output, err := transform(dir, sugar.Config{}, mtc)
	if err != nil {
		return nil, err
	}

	fm := genlib.NewFileManager(dir)
	if err = fm.Add(&file.GoFile{Path: filepath.Join(dir, outFileName), Content: string(output)}); err != nil {
		return nil, fmt.Errorf("cannot add output.go to FileManager: %w", err)
	}
	return fm, nil
}
