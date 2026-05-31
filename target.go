package sugar

import (
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
)

type FilePath struct {
	RelPath     string
	AbsPath     string
	DisplayPath string
}

type TargetCollection []Target

func (c TargetCollection) Resolve() ([]FilePath, error) {
	var result []FilePath
	for _, v := range c {
		rs, err := v.Resolve()
		if err != nil {
			return nil, err
		}
		for _, r := range rs {
			result = append(result, r)
		}
	}
	slices.SortFunc(result, func(a, b FilePath) int {
		return strings.Compare(a.DisplayPath, b.DisplayPath)
	})
	return result, nil
}

type Target struct {
	Root       string
	WorkingDir string
	Input      string
	Config     Config

	Path      string
	IsDir     bool
	Recursive bool
}

func (t *Target) DisplayPath() string {
	displayPath, err := filepath.Rel(t.WorkingDir, t.Path)
	if err != nil {
		displayPath = t.Path
	}
	return displayPath
}

func (t *Target) Resolve() ([]FilePath, error) {
	var result []FilePath
	if t.IsDir {
		if err := t.walkDir(t.Path, &result, t.Recursive); err != nil {
			return nil, err
		}
		return result, nil
	}

	// resolve the Target as file
	fp, err := t.makeFilePath(t.Path)
	if err != nil {
		return nil, err
	}
	if fp != nil {
		result = append(result, *fp)
	}
	return result, nil
}

func (t *Target) Dirs() ([]string, error) {
	var result []string
	err := filepath.WalkDir(t.Path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			return nil
		}

		if !t.Recursive && path != t.Path {
			return fs.SkipDir
		}
		result = append(result, path)
		return nil
	})
	return result, err
}

func (t *Target) IsSugarFilePath(path string) (FilePath, bool) {
	f, err := t.makeFilePath(path)
	if f == nil || err != nil {
		return FilePath{}, false
	}
	return *f, true
}

func (t *Target) walkDir(dir string, result *[]FilePath, recursive bool) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if !recursive && path != dir {
				return fs.SkipDir
			}
			return nil
		}

		fp, err := t.makeFilePath(path)
		if err != nil {
			return err
		}

		if fp != nil {
			*result = append(*result, *fp)
		}
		return nil
	})
}

func (t *Target) makeFilePath(path string) (*FilePath, error) {
	ext := filepath.Ext(path)
	if !t.Config.Env.IsSugarFile(ext) {
		return nil, nil
	}

	relPath, err := filepath.Rel(t.Root, path)
	if err != nil {
		return nil, err
	}
	displayPath, err := filepath.Rel(t.WorkingDir, path)
	if err != nil {
		return nil, err
	}
	return &FilePath{RelPath: relPath, AbsPath: path, DisplayPath: displayPath}, nil
}
