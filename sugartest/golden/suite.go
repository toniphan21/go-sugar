package golden

import (
	"embed"
	"io/fs"
	"testing"

	"nhatp.com/go/gen-lib/gentest"
	"nhatp.com/go/sugar/sugartest"
)

type TestCase = gentest.MarkdownTestCase

type TestSuite struct {
	Name  string
	FS    embed.FS
	Run   func(t *testing.T, tc TestCase)
	Match string
	File  string
}

func (t *TestSuite) FilePaths() []string {
	var result []string
	switch {
	case t.Match != "":
		files, err := fs.Glob(t.FS, t.Match)
		if err == nil {
			return files
		}

	case t.File != "":
		f, err := t.FS.Open(t.File)
		if err == nil {
			stat, err := f.Stat()
			if err == nil && !stat.IsDir() {
				result = append(result, t.File)
			}
		}
	}
	return result
}

func Test(t *testing.T, fs embed.FS, file string, fn func(t *testing.T, tc TestCase)) {
	gentest.RunEmbedGoldenFile(t, fs, file, fn)
}

func GeneratePipeline(t *testing.T, tc TestCase) {
	t.Helper()
	sugartest.RunGoldenGeneratePipelineTest(t, tc)
}

func FormatPipeline(t *testing.T, tc TestCase) {
	t.Helper()
	sugartest.RunGoldenFormatPipelineTest(t, tc)
}

func RestoreTransform(t *testing.T, tc TestCase) {
	t.Helper()
	sugartest.RunGoldenRestoreTransformTest(t, tc)
}

func SemanticTransform(t *testing.T, tc TestCase) {
	t.Helper()
	sugartest.RunGoldenSemanticTransformTest(t, tc)
}

func StructuralTransform(t *testing.T, tc TestCase) {
	t.Helper()
	sugartest.RunGoldenStructuralTransformTest(t, tc)
}
