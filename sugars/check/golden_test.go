package check

import (
	"embed"
	"testing"

	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/sugartest/golden"
)

//go:embed testdata
var testdata embed.FS

func Test_Golden(t *testing.T) {
	sugar.Register(New())

	suite := []golden.TestSuite{
		{Name: "generate", FS: testdata, Match: "testdata/generate-pipeline/*.md", Run: golden.GeneratePipeline},
		{Name: "format", FS: testdata, Match: "testdata/format-pipeline/*.md", Run: golden.FormatPipeline},
		{Name: "t1", FS: testdata, Match: "testdata/structural-transform/*.md", Run: golden.StructuralTransform},
		{Name: "t2", FS: testdata, Match: "testdata/semantic-transform/*.md", Run: golden.SemanticTransform},
		{Name: "t3", FS: testdata, Match: "testdata/restore-transform/*.md", Run: golden.RestoreTransform},

		{
			Name: "generate-dev", FS: testdata,
			File: "testdata/generate-pipeline/basic.md",
			Run:  golden.GeneratePipeline,
		},

		{
			Name: "format-dev", FS: testdata,
			File: "testdata/format-pipeline/basic.md",
			Run:  golden.FormatPipeline,
		},
	}

	for _, tt := range suite {
		t.Run(tt.Name, func(t *testing.T) {
			for _, p := range tt.FilePaths() {
				golden.Test(t, tt.FS, p, func(t *testing.T, tc golden.TestCase) {
					tt.Run(t, tc)
				})
			}
		})
	}
}
