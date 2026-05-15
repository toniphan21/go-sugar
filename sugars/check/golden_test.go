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
	gentest.RunEmbedGoldenFiles(t, goldenStructuralTransformMarkdownFiles, func(t *testing.T, testCase gentest.MarkdownTestCase) {
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
			gentest.RunEmbedGoldenFile(t, goldenStructuralTransformMarkdownFiles, tc.file, func(t *testing.T, testCase gentest.MarkdownTestCase) {
				runGoldenStructuralTransformTest(t, testCase)
			})
		})
	}
}

func runGoldenStructuralTransformTest(t *testing.T, testCase gentest.MarkdownTestCase) {
	t.Helper()

	sugar.Register(New())

	sugartest.RunGoldenStructuralTransformTest(t, testCase)
}

// ---

//go:embed testdata/semantic-transform/*.md
var goldenSemanticTransformMarkdownFiles embed.FS

func TestGoldenSemanticTransformFiles(t *testing.T) {
	gentest.RunEmbedGoldenFiles(t, goldenSemanticTransformMarkdownFiles, func(t *testing.T, testCase gentest.MarkdownTestCase) {
		runGoldenSemanticTransformTest(t, testCase)
	})
}

func TestGoldenSemanticTransformFiles_Dev(t *testing.T) {
	cases := []struct {
		file string
	}{
		{file: "testdata/semantic-transform/basic.md"},
	}

	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			gentest.RunEmbedGoldenFile(t, goldenSemanticTransformMarkdownFiles, tc.file, func(t *testing.T, testCase gentest.MarkdownTestCase) {
				runGoldenSemanticTransformTest(t, testCase)
			})
		})
	}
}

func runGoldenSemanticTransformTest(t *testing.T, testCase gentest.MarkdownTestCase) {
	t.Helper()

	sugar.Register(New())

	sugartest.RunGoldenSemanticTransformTest(t, testCase)
}
