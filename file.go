package sugar

import (
	"fmt"
	"slices"
	"strings"
)

func newFile(relPath, goFilePath string, content []byte) *File {
	snapshot := newSnapshot(content)
	return &File{
		sugarPath: relPath,
		goPath:    goFilePath,
		current:   snapshot,
	}
}

type FileError struct {
	File   File
	Action string
	Err    error
}

func (e *FileError) Unwrap() error {
	return e.Err
}

func (e *FileError) Error() string {
	return fmt.Sprintf("%v: %v: %v", e.File.sugarPath, e.Action, e.Err)
}

type FileErrors []*FileError

func (fe FileErrors) GroupByPath() map[string][]*FileError {
	grouped := make(map[string][]*FileError)
	for _, v := range fe {
		grouped[v.File.sugarPath] = append(grouped[v.File.sugarPath], v)
	}
	return grouped
}

func (fe FileErrors) Unwrap() []error {
	errs := make([]error, len(fe))
	for i, e := range fe {
		errs[i] = e
	}
	return errs
}

func (fe FileErrors) Error() string {
	sb := strings.Builder{}
	paths := make([]string, len(fe))
	grouped := fe.GroupByPath()
	for v := range grouped {
		paths = append(paths, v)
	}

	slices.Sort(paths)
	for _, path := range paths {
		errs := grouped[path]
		for _, err := range errs {
			sb.WriteString(err.Error() + "\n")
		}
	}
	return sb.String()
}

type File struct {
	hash      [32]byte
	sugarPath string
	goPath    string
	current   *Snapshot
}

func (f *File) SugarPath() string {
	return f.sugarPath
}

func (f *File) GoPath() string {
	return f.goPath
}

func (f *File) Update(content []byte) {
	f.current = newSnapshot(content)
}

func (f *File) Content() []byte {
	return f.current.source
}

func (f *File) StructuralTransform() []byte {
	return f.current.StructuralTransform()
}

func (f *File) SemanticTransform(module ModuleScope) ([]byte, error) {
	t1smap := f.current.t1smap
	if t1smap == nil {
		return nil, fmt.Errorf("missing t1smap: call StructuralTransform first")
	}

	return f.current.SemanticTransform(module, FileScope{
		PkgPath:     module.ResolvePackagePath(f.sugarPath),
		T1SourceMap: *t1smap,
	})
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
