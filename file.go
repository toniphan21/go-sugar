package sugar

import (
	"golang.org/x/tools/go/packages"
)

func newFile(relPath, goFilePath string, content []byte) *File {
	snapshot := newSnapshot(content)
	return &File{
		sugarPath: relPath,
		goPath:    goFilePath,
		current:   snapshot,
	}
}

type File struct {
	hash      [32]byte
	content   []byte
	sugarPath string
	goPath    string
	current   *Snapshot
}

func (f *File) Update(content []byte) {
	f.current = newSnapshot(content)
}

func (f *File) structuralTransform() error {
	return f.current.structuralTransform()
}

func (f *File) semanticAnalysis(pkg *packages.Package) error {
	return f.current.semanticAnalysis(pkg)
}

func (f *File) StructuralTransform() []byte {
	return f.current.StructuralTransform()
}

func (f *File) semanticTransform() error {
	return f.current.semanticTransform()
}

func (f *File) SemanticTransform() []byte {
	return f.current.SemanticTransform()
}

func (f *File) Hash() [32]byte {
	return f.current.Hash()
}

func (f *File) SugarToGo(line, column int) (int, int, error) {
	return f.current.SugarToGo(line, column)
}

func (f *File) GoToSugar(line, column int) (int, int, error) {
	return f.current.GoToSugar(line, column)
}

var _ fileAPI = (*File)(nil)
