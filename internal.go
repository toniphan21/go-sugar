package sugar

type shareAPI interface {
	Hash() [32]byte

	StructuralTransform() []byte

	SugarToGo(line, column int) (int, int, error)

	GoToSugar(line, column int) (int, int, error)
}

type snapshotAPI interface {
	shareAPI

	SemanticTransform(module ModuleScope, file FileScope) ([]byte, error)
}

type fileAPI interface {
	shareAPI

	SemanticTransform(module ModuleScope) ([]byte, error)

	SugarPath() string

	GoPath() string

	Update(content []byte)

	Content() []byte
}

type moduleAPI interface {
	RegisterFile(relPath string) (*File, error)

	AddFile(relPath string, content []byte) (*File, error)

	RemoveFile(relPath string)

	File(relPath string) (*File, bool)

	HasFile(relPath string) bool
}
