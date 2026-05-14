package sugar

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"
)

var ErrCannotPerformActionOnGoFile = errors.New("cannot perform action on go file")

func NewModule(path string, config Config) (*Module, error) {
	env := config.env()
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	var dir string
	if info.IsDir() {
		dir = path
	} else {
		dir = filepath.Dir(path)
	}

	// walk up to find go.mod
	root := dir
	for {
		goMod := filepath.Join(root, env.GoModFileName)
		if _, err := os.Stat(goMod); err == nil {
			break
		}
		parent := filepath.Dir(root)
		if parent == root {
			return nil, errors.New(env.GoModFileName + " not found")
		}
		root = parent
	}

	// read module path from go.mod
	goModPath := filepath.Join(root, env.GoModFileName)
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return nil, err
	}
	pkgPath := modfile.ModulePath(data)
	if pkgPath == "" {
		return nil, errors.New("module path not found in " + env.GoModFileName)
	}

	mod := &Module{
		Root:        root,
		GoModPath:   goModPath,
		PackagePath: pkgPath,
		Config:      config,
		files:       make(map[string]*File),
	}
	if err = mod.DiscoverFiles(); err != nil {
		return nil, err
	}
	return mod, nil
}

// ---

type Module struct {
	Root        string
	GoModPath   string
	PackagePath string
	Config      Config
	files       map[string]*File
}

func (m *Module) DiscoverFiles() error {
	env := m.Config.env()
	return filepath.WalkDir(m.Root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if !env.IsGoFile(ext) && !env.IsSugarFile(ext) {
			return nil
		}

		relPath, err := filepath.Rel(m.Root, path)
		if err != nil {
			return err
		}

		_, err = m.RegisterFile(relPath)
		return err
	})
}

func (m *Module) RegisterFile(relPath string) (*File, error) {
	absPath := filepath.Join(m.Root, relPath)
	content, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	return m.AddFile(relPath, content)
}

func (m *Module) AddFile(relPath string, content []byte) (*File, error) {
	env := m.Config.env()
	ext := filepath.Ext(relPath)

	if !env.IsGoFile(ext) && !env.IsSugarFile(ext) {
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}

	if env.IsSugarFile(ext) {
		file := newSugarFile(relPath, env.GoFilePath(relPath), content)
		m.files[relPath] = file
		return file, nil
	}

	file := newGoFile(relPath, content)
	m.files[relPath] = file
	return file, nil
}

func (m *Module) RemoveFile(relPath string) {
	delete(m.files, relPath)
}

func (m *Module) File(relPath string) (*File, bool) {
	f, ok := m.files[relPath]
	return f, ok
}

func (m *Module) HasFile(relPath string) bool {
	_, ok := m.files[relPath]
	return ok
}

func (m *Module) Transform() error {
	err := m.structuralTransform()
	if err != nil {
		return err
	}

	var sugarFiles = make(map[string]*File)
	overlay := make(map[string][]byte)
	for p, f := range m.files {
		if f.isSugar {
			sugarFiles[p] = f
			fp := filepath.Join(m.Root, f.goPath)
			overlay[fp] = f.StructuralTransform()
		}
	}

	cfg := &packages.Config{
		Mode:    packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedDeps | packages.NeedImports,
		Fset:    token.NewFileSet(),
		Overlay: overlay,
		Dir:     m.Root,
	}
	pkgs, _ := packages.Load(cfg, "./...")
	for p, f := range sugarFiles {
		fpp := m.filePkgPath(p)

		var pkg *packages.Package
		for _, v := range pkgs {
			if fpp == v.ID {
				pkg = v
				break
			}
		}

		if pkg != nil {
			// not fail just because the semantic analysis fail
			_ = f.semanticAnalysis(pkg)
		}
	}
	return nil
}

func (m *Module) structuralTransform() error {
	for _, v := range m.files {
		if !v.isSugar {
			continue
		}

		if err := v.structuralTransform(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Module) filePkgPath(relPath string) string {
	dir := filepath.Dir(relPath)
	if dir == "." {
		return m.PackagePath
	}
	return m.PackagePath + "/" + filepath.ToSlash(dir)
}

var _ moduleAPI = (*Module)(nil)

// ---

func newSugarFile(relPath, goFilePath string, content []byte) *File {
	snapshot := newSnapshot(content)
	return &File{
		isSugar:   true,
		sugarPath: relPath,
		goPath:    goFilePath,
		snapshot:  snapshot,
	}
}

func newGoFile(relPath string, content []byte) *File {
	return &File{
		isSugar: false,
		goPath:  relPath,
		hash:    sha256.Sum256(content),
		content: content,
	}
}

type File struct {
	isSugar   bool
	hash      [32]byte
	content   []byte
	sugarPath string
	goPath    string
	snapshot  *Snapshot
}

func (f *File) Update(content []byte) {
	if f.isSugar {
		f.snapshot = newSnapshot(content)
		return
	}

	f.hash = sha256.Sum256(content)
	f.content = content
}

func (f *File) structuralTransform() error {
	if f.isSugar {
		return f.snapshot.structuralTransform()
	}
	return ErrCannotPerformActionOnGoFile
}

func (f *File) semanticAnalysis(pkg *packages.Package) error {
	if f.isSugar {
		return f.snapshot.semanticAnalysis(pkg)
	}
	return ErrCannotPerformActionOnGoFile
}

func (f *File) StructuralTransform() []byte {
	if f.isSugar {
		return f.snapshot.StructuralTransform()
	}
	return nil
}

func (f *File) Transform() []byte {
	if f.isSugar {
		return f.snapshot.Transform()
	}
	return nil
}

func (f *File) Hash() [32]byte {
	if f.isSugar {
		return f.snapshot.Hash()
	}
	return f.hash
}

func (f *File) SugarToGo(line, column int) (int, int, error) {
	if f.isSugar {
		return f.snapshot.SugarToGo(line, column)
	}
	return line, column, nil
}

func (f *File) GoToSugar(line, column int) (int, int, error) {
	if f.isSugar {
		return f.snapshot.GoToSugar(line, column)
	}
	return line, column, nil
}

var _ fileAPI = (*File)(nil)
