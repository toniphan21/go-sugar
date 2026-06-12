package sugartest

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	genlib "nhatp.com/go/gen-lib"
	"nhatp.com/go/gen-lib/cli"
	"nhatp.com/go/gen-lib/file"
	"nhatp.com/go/gen-lib/gentest"
	"nhatp.com/go/sugar"
)

func goldenFile(tc gentest.MarkdownTestCase, filename string) ([]byte, error) {
	var output []byte
	for _, v := range tc.GoldenFiles {
		if v.FilePath() == filename {
			output = v.FileContent()
		}
	}
	if output == nil {
		return nil, fmt.Errorf("no there is no expected output, use `// golden-file: %v` in a codeblock", filename)
	}
	return output, nil
}

func makeModuleFromMarkdownTestCase(dir string, config sugar.Config, tc gentest.MarkdownTestCase) (*sugar.Module, error) {
	if err := genlib.SetupSourceCode(dir, tc.SourceFiles); err != nil {
		return nil, err
	}

	return sugar.NewModule(dir, config)
}

func assertFileOutput(t *testing.T, actual, expected []byte) {
	t.Helper()

	if !bytes.Equal(actual, expected) {
		if !PrintColor {
			cli.DisableColor()
		}

		cli.PrintDiffWithFunction("expected", expected, "actual", actual, func(line string) {
			if PrintDiff {
				fmt.Println(line)
			}
		})
		t.Fatal("output does not match the golden file")
	}
}

func RunGoldenGeneratePipelineTest(t *testing.T, tc gentest.MarkdownTestCase) {
	t.Helper()
	dir := t.TempDir()

	fm, err := PerformGeneratePipeline(dir, sugar.DefaultConfig(), tc)
	if err != nil {
		t.Fatal(err)
	}

	var out = make(map[string][]byte) // path -> content
	for _, f := range fm.Files() {
		absPath := f.FilePath()
		relPath, err := filepath.Rel(dir, absPath)
		if err != nil {
			t.Fatal(err)
		}
		out[relPath] = f.FileContent()
	}

	for _, f := range tc.GoldenFiles {
		result, ok := out[f.FilePath()]
		if !ok {
			t.Errorf(`expected file "%s" but file is not generated`, f.FilePath())
		}

		assertFileOutput(t, result, f.FileContent())
	}
}

func PerformGeneratePipeline(dir string, config sugar.Config, tc gentest.MarkdownTestCase) (genlib.FileManager, error) {
	if err := genlib.SetupSourceCode(dir, tc.SourceFiles); err != nil {
		return nil, err
	}

	mod, err := sugar.NewModule(dir, sugar.DefaultConfig())
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

func RunGoldenFormatPipelineTest(t *testing.T, tc gentest.MarkdownTestCase) {
	t.Helper()

	result, err := PerformFormatPipeline(t.TempDir(), sugar.DefaultConfig(), tc)
	if err != nil {
		t.Fatal(err)
	}

	output, err := goldenFile(tc, "output.gos")
	if err != nil {
		t.Fatal(err.Error())
	}

	assertFileOutput(t, result, output)
}

func PerformFormatPipeline(dir string, config sugar.Config, tc gentest.MarkdownTestCase) ([]byte, error) {
	mod, err := makeModuleFromMarkdownTestCase(dir, config, tc)
	if err != nil {
		return nil, err
	}

	f, ok := mod.File("input.gos")
	if !ok {
		return nil, errors.New("cannot find input.gos file")
	}
	return mod.FormatFile(f)
}

func RunGoldenStructuralTransformTest(t *testing.T, tc gentest.MarkdownTestCase) {
	t.Helper()

	result, err := PerformStructuralTransform(t.TempDir(), sugar.DefaultConfig(), tc)
	if err != nil {
		t.Fatal(err)
	}

	output, err := goldenFile(tc, "output.go")
	if err != nil {
		t.Fatal(err.Error())
	}

	assertFileOutput(t, result, output)
}

func PerformStructuralTransform(dir string, config sugar.Config, tc gentest.MarkdownTestCase) ([]byte, error) {
	mod, err := makeModuleFromMarkdownTestCase(dir, config, tc)
	if err != nil {
		return nil, err
	}

	mod.StructuralTransform()
	f, ok := mod.File("input.gos")
	if !ok {
		return nil, errors.New("cannot find input.gos file")
	}
	return f.StructuralTransform(), nil
}

func RunGoldenSemanticTransformTest(t *testing.T, tc gentest.MarkdownTestCase) {
	t.Helper()

	result, err := PerformSemanticTransform(t.TempDir(), sugar.DefaultConfig(), tc)
	if err != nil {
		t.Fatal(err)
	}

	output, err := goldenFile(tc, "output.go")
	if err != nil {
		t.Fatal(err.Error())
	}

	assertFileOutput(t, result, output)
}

func PerformSemanticTransform(dir string, config sugar.Config, tc gentest.MarkdownTestCase) ([]byte, error) {
	mod, err := makeModuleFromMarkdownTestCase(dir, config, tc)
	if err != nil {
		return nil, err
	}

	if err = mod.SemanticTransform(); err != nil {
		return nil, err
	}

	f, ok := mod.File("input.gos")
	if !ok {
		return nil, errors.New("cannot find input.gos file")
	}
	return f.SemanticTransform(mod.Scope())
}

func RunGoldenRestoreTransformTest(t *testing.T, tc gentest.MarkdownTestCase) {
	t.Helper()

	result, err := PerformRestoreTransform(t.TempDir(), sugar.DefaultConfig(), tc)
	if err != nil {
		t.Fatal(err)
	}

	output, err := goldenFile(tc, "output.gos")
	if err != nil {
		t.Fatal(err.Error())
	}

	assertFileOutput(t, result, output)
}

func PerformRestoreTransform(_ string, _ sugar.Config, tc gentest.MarkdownTestCase) ([]byte, error) {
	var input []byte
	for _, v := range tc.SourceFiles {
		if v.FilePath() == "input.go" {
			input = v.FileContent()
		}
	}
	if input == nil {
		return nil, errors.New("no there is no input, use `// file: input.go` in a codeblock")
	}

	return sugar.RestoreTransform(input), nil
}
