package sugar

import "golang.org/x/tools/go/packages"

type snapshotAPI interface {
	Hash() [32]byte

	structuralTransform() error

	StructuralTransform() []byte

	semanticAnalysis(pkg *packages.Package) error

	semanticTransform() error

	SemanticTransform() []byte

	SugarToGo(line, column int) (int, int, error)

	GoToSugar(line, column int) (int, int, error)
}

type fileAPI interface {
	snapshotAPI

	Update(content []byte)
}

type moduleAPI interface {
	RegisterFile(relPath string) (*File, error)

	AddFile(relPath string, content []byte) (*File, error)

	RemoveFile(relPath string)

	File(relPath string) (*File, bool)

	HasFile(relPath string) bool
}
