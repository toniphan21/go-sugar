package check

import (
	"embed"
	"testing"

	"nhatp.com/go/gen-lib/gentest"
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/sugartest"
)

//go:embed testdata/structural-transform/*.md
var mdStructuralTransform embed.FS

func Test_GoldenStructuralTransform_Files(t *testing.T) {
	fs := mdStructuralTransform
	run := func(t *testing.T, testCase gentest.MarkdownTestCase) {
		sugar.Register(New())
		sugartest.RunGoldenStructuralTransformTest(t, testCase)
	}

	t.Run("dev", func(t *testing.T) {
		file := "testdata/structural-transform/basic.md"
		gentest.RunEmbedGoldenFile(t, fs, file, run)
	})

	t.Run("embed", func(t *testing.T) {
		gentest.RunEmbedGoldenFiles(t, fs, run)
	})
}

//go:embed testdata/semantic-transform/*.md
var mdSemanticTransform embed.FS

func Test_GoldenSemanticTransform_Files(t *testing.T) {
	fs := mdSemanticTransform
	run := func(t *testing.T, testCase gentest.MarkdownTestCase) {
		sugar.Register(New())
		sugartest.RunGoldenSemanticTransformTest(t, testCase)
	}

	t.Run("dev", func(t *testing.T) {
		file := "testdata/semantic-transform/basic.md"
		gentest.RunEmbedGoldenFile(t, fs, file, run)
	})

	t.Run("embed", func(t *testing.T) {
		gentest.RunEmbedGoldenFiles(t, fs, run)
	})
}

//go:embed testdata/restore-transform/*.md
var mdRestoreTransform embed.FS

func Test_GoldenRestoreTransform_Files(t *testing.T) {
	fs := mdRestoreTransform
	run := func(t *testing.T, testCase gentest.MarkdownTestCase) {
		sugar.Register(New())
		sugartest.RunGoldenRestoreTransformTest(t, testCase)
	}

	t.Run("dev", func(t *testing.T) {
		file := "testdata/restore-transform/basic.md"
		gentest.RunEmbedGoldenFile(t, fs, file, run)
	})

	t.Run("embed", func(t *testing.T) {
		gentest.RunEmbedGoldenFiles(t, fs, run)
	})
}
