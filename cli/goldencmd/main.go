package goldencmd

import (
	"fmt"
	"io"
	"log/slog"
	"os"
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
	TypeGeneratePipeline Type = iota
	TypeStructuralTransform
	TypeSemanticTransform
	TypeRestoreTransform
	TypeFormatPipeline
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

func Run(stdin, stdout, stderr *os.File, args Arguments) error {
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
	case TypeGeneratePipeline:
		runGenerate(cmd)
	case TypeStructuralTransform:
		runStructuralTransform(cmd)
	case TypeSemanticTransform:
		runSemanticTransform(cmd)
	case TypeRestoreTransform:
		runRestoreTransform(cmd)
	case TypeFormatPipeline:
		runFormat(cmd)
	}
	return nil
}

func runGenerate(cmd *cli.TestRunner) {
	printTestType(cmd, "Generate Pipeline")
	cmd.RunTestCase = func(tc cli.TestCase, options map[string]any) (genlib.FileManager, error) {
		mtc := gentest.MarkdownTestCase{
			SourceFiles: tc.SourceFiles,
		}
		return sugartest.PerformGeneratePipeline(tc.TestDir, sugar.DefaultConfig(), mtc)
	}
	cmd.Run()
}

func runStructuralTransform(cmd *cli.TestRunner) {
	printTestType(cmd, "T1 - StructuralTransform")
	cmd.RunTestCase = func(tc cli.TestCase, options map[string]any) (genlib.FileManager, error) {
		return runTransformTest(tc, "output.go", sugartest.PerformStructuralTransform)
	}
	cmd.Run()
}

func runSemanticTransform(cmd *cli.TestRunner) {
	printTestType(cmd, "T2 - SemanticTransform")
	cmd.RunTestCase = func(tc cli.TestCase, options map[string]any) (genlib.FileManager, error) {
		return runTransformTest(tc, "output.go", sugartest.PerformSemanticTransform)
	}
	cmd.Run()
}

func runRestoreTransform(cmd *cli.TestRunner) {
	printTestType(cmd, "T3 - RestoreTransform")
	cmd.RunTestCase = func(tc cli.TestCase, options map[string]any) (genlib.FileManager, error) {
		return runTransformTest(tc, "output.gos", sugartest.PerformRestoreTransform)
	}
	cmd.Run()
}

func runFormat(cmd *cli.TestRunner) {
	printTestType(cmd, "Format Pipeline")
	cmd.RunTestCase = func(tc cli.TestCase, options map[string]any) (genlib.FileManager, error) {
		return runTransformTest(tc, "output.gos", sugartest.PerformFormatPipeline)
	}
	cmd.Run()
}

func runTransformTest(tc cli.TestCase, outFileName string, transform func(string, sugar.Config, gentest.MarkdownTestCase) ([]byte, error)) (genlib.FileManager, error) {
	dir := tc.TestDir
	mtc := gentest.MarkdownTestCase{
		SourceFiles: tc.SourceFiles,
	}

	output, err := transform(dir, sugar.DefaultConfig(), mtc)
	if err != nil {
		return nil, err
	}

	fm := genlib.NewFileManager(dir)
	if err = fm.Add(&file.GoFile{Path: filepath.Join(dir, outFileName), Content: string(output)}); err != nil {
		return nil, fmt.Errorf("cannot add output.go to FileManager: %w", err)
	}
	return fm, nil
}

func printTestType(cmd *cli.TestRunner, typeName string) {
	cmd.Print("Running " + cli.ColorYellow(typeName) + " golden tests with " + color.Binary(sugar.BinaryName) + " " + color.Version(sugar.BinaryVersion))
	cmd.Print("")
}
