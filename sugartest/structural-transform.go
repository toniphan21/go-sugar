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

//type GoldenStructuralTransformTestCase struct {
//	Input     string
//	Output    string
//	SourceMap string
//}

func RunGoldenStructuralTransformTest(t *testing.T, tc gentest.MarkdownTestCase) {
	t.Helper()

	var input, output []byte
	for _, v := range tc.SourceFiles {
		if v.FilePath() == "input.gos" {
			input = v.FileContent()
		}
	}
	if input == nil {
		t.Fatal("no there is no input, use `// file: input.gos` in a codeblock")
	}

	for _, v := range tc.GoldenFiles {
		if v.FilePath() == "output.go" {
			output = v.FileContent()
		}
	}

	if output == nil {
		t.Fatal("no there is no expected output, use `// golden-file: output.go` in a codeblock")
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

	err = mod.Transform()
	if err != nil {
		t.Fatal(err.Error())
	}

	f, ok := mod.File("input.gos")
	if !ok {
		t.Fatal("cannot find input.gos file")
	}

	result := f.StructuralTransform()
	if !bytes.Equal(result, output) {
		if !PrintColor {
			cli.DisableColor()
		}

		cli.PrintDiffWithFunction("expected", output, "actual", result, func(line string) {
			if PrintDiff {
				fmt.Println(line)
			}
		})
		t.Fatal("output does not match the golden file")
	}
}
