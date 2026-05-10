package check

import (
	"embed"
	"testing"

	"nhatp.com/go/gen-lib/gentest"
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/sugartest"
)

//go:embed testdata/structural-transform/*.md
var goldenStructuralTransformMarkdownFiles embed.FS

func TestGoldenStructuralTransformFiles(t *testing.T) {
	gentest.RunEmbedGoldenFiles(t, goldenStructuralTransformMarkdownFiles, func(testCase gentest.MarkdownTestCase) {
		runGoldenStructuralTransformTest(t, testCase)
	})
}

func TestGoldenStructuralTransformFiles_Dev(t *testing.T) {
	cases := []struct {
		file string
	}{
		{file: "testdata/structural-transform/basic.md"},
	}

	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			gentest.RunEmbedGoldenFile(t, goldenStructuralTransformMarkdownFiles, tc.file, func(testCase gentest.MarkdownTestCase) {
				runGoldenStructuralTransformTest(t, testCase)
			})
		})
	}
}

func runGoldenStructuralTransformTest(t *testing.T, testCase gentest.MarkdownTestCase) {
	sugar.Register(New())

	sugartest.RunGoldenStructuralTransformTest(t, testCase)
}
