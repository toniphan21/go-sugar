package sugartest

import (
	"bytes"
	"fmt"
	"testing"

	genlib "nhatp.com/go/gen-lib"
	"nhatp.com/go/gen-lib/cli"
	"nhatp.com/go/gen-lib/gentest"
	"nhatp.com/go/sugar"
)

func goldenFile(t *testing.T, tc gentest.MarkdownTestCase, filename string) []byte {
	t.Helper()

	var output []byte
	for _, v := range tc.GoldenFiles {
		if v.FilePath() == filename {
			output = v.FileContent()
		}
	}
	if output == nil {
		t.Fatalf("no there is no expected output, use `// golden-file: %d` in a codeblock", filename)
	}
	return output
}

func makeModuleFromMarkdownTestCase(t *testing.T, tc gentest.MarkdownTestCase) *sugar.Module {
	t.Helper()

	var input []byte
	for _, v := range tc.SourceFiles {
		if v.FilePath() == "input.gos" {
			input = v.FileContent()
		}
	}
	if input == nil {
		t.Fatal("no there is no input, use `// file: input.gos` in a codeblock")
	}

	dir := t.TempDir()
	err := genlib.SetupSourceCode(dir, tc.SourceFiles)
	if err != nil {
		t.Fatal(err.Error())
	}

	mod, err := sugar.NewModule(dir, sugar.Config{})
	if err != nil {
		t.Fatal(err.Error())
	}
	return mod
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

func RunGoldenStructuralTransformTest(t *testing.T, tc gentest.MarkdownTestCase) {
	t.Helper()

	mod := makeModuleFromMarkdownTestCase(t, tc)
	if err := mod.StructuralTransform(); err != nil {
		t.Fatal(err.Error())
	}

	f, ok := mod.File("input.gos")
	if !ok {
		t.Fatal("cannot find input.gos file")
	}

	output := goldenFile(t, tc, "output.go")
	result := f.StructuralTransform()

	assertFileOutput(t, result, output)
}

func RunGoldenSemanticTransformTest(t *testing.T, tc gentest.MarkdownTestCase) {
	t.Helper()

	mod := makeModuleFromMarkdownTestCase(t, tc)
	if err := mod.SemanticTransform(); err != nil {
		t.Fatal(err.Error())
	}

	f, ok := mod.File("input.gos")
	if !ok {
		t.Fatal("cannot find input.gos file")
	}

	output := goldenFile(t, tc, "output.go")
	result := f.SemanticTransform()

	assertFileOutput(t, result, output)
}
