package sugar

import (
	"errors"
	"fmt"
	"go/token"
	"path/filepath"

	"github.com/oklog/ulid/v2"
	"golang.org/x/tools/go/packages"
)

type SemanticScope struct {
	ModuleScope `json:"module"`
	FileScope   `json:"file"`
}

type ModuleScope struct {
	Overlay    map[string][]byte `json:"overlay"`
	Root       string            `json:"root"`
	ModulePath string            `json:"modulePath"`
}

func (m *ModuleScope) ResolvePackagePath(relPath string) string {
	dir := filepath.Dir(relPath)
	if dir == "." {
		return m.ModulePath
	}
	return m.ModulePath + "/" + filepath.ToSlash(dir)
}

type FileScope struct {
	PkgPath     string    `json:"pkgPath"`
	T1SourceMap SourceMap `json:"t1SourceMap"`
}

type BinaryInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Usage   string `json:"usage"`
}

type Sugar interface {
	ID() string

	Binary() BinaryInfo

	Parse(source []byte) []Node

	Restore(source []byte) []Node

	PrepareSource(id string, source []byte)

	PrepareSemanticScope(id string, scope SemanticScope)

	CleanUp(sourceID string, scopeID string)

	StructuralTransform(sourceID string, n Node) ([]byte, error)

	SemanticTransformer(sourceID string, scopeID string, n Node) ([]byte, error)

	RestoreTransform(sourceID string, n Node) ([]byte, error)
}

// --- internal

var ErrUnknownNode = errors.New("unknown node")
var ErrUnpreparedSource = errors.New("unprepare source")
var ErrUnpreparedSemanticScope = errors.New("unprepare semantic scope")
var ErrUnknownPkg = errors.New("unknown pkg")

type parsedNode struct {
	sugar Sugar
	node  Node
}

func makeSourceID(source []byte) string {
	return ulid.Make().String()
}

func makeScopeID(_ ModuleScope, _ FileScope) string {
	return ulid.Make().String()
}

// --- helpers

type Base struct {
	Sources  map[string][]byte
	Scopes   map[string]SemanticScope
	Packages map[string]*packages.Package
}

func (b *Base) PrepareSource(id string, source []byte) {
	if b.Sources == nil {
		b.Sources = make(map[string][]byte)
	}
	b.Sources[id] = source
}

func (b *Base) PrepareSemanticScope(id string, scope SemanticScope) {
	if b.Scopes == nil {
		b.Scopes = make(map[string]SemanticScope)
	}
	b.Scopes[id] = scope
}

func (b *Base) CleanUp(sourceID, scopeID string) {
	if b.Sources != nil {
		delete(b.Sources, sourceID)
	}
	if b.Scopes != nil {
		delete(b.Scopes, scopeID)
	}
}

func DoTransform[T any](b *Base, sourceID string, n Node, fn func(source []byte, data T) ([]byte, error)) ([]byte, error) {
	data, ok := n.(T)
	if !ok {
		return nil, fmt.Errorf("%w: expected T, got %T", ErrUnknownNode, n)
	}

	source, have := b.Sources[sourceID]
	if !have {
		return nil, ErrUnpreparedSource
	}
	return fn(source, data)
}

func DoTransformWithSemanticScope[T any](b *Base, sourceID, scopeID string, n Node, fn func(source []byte, data T) ([]byte, error)) ([]byte, error) {
	data, ok := n.(T)
	if !ok {
		return nil, ErrUnknownNode
	}

	source, have := b.Sources[sourceID]
	if !have {
		return nil, ErrUnpreparedSource
	}

	scope, have := b.Scopes[scopeID]
	if !have {
		return nil, ErrUnpreparedSemanticScope
	}

	sn, ok := n.(SemanticNode)
	if !ok {
		return fn(source, data)
	}

	if b.Packages == nil {
		b.Packages = make(map[string]*packages.Package)
	}
	cached, have := b.Packages[scopeID]
	if !have {
		cfg := &packages.Config{
			Mode:    packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedDeps | packages.NeedImports,
			Fset:    token.NewFileSet(),
			Overlay: scope.Overlay,
			Dir:     scope.Root,
		}
		pkgs, _ := packages.Load(cfg, "./...")

		var pkg *packages.Package
		for _, p := range pkgs {
			if scope.PkgPath == p.ID {
				pkg = p
				break
			}
		}

		if pkg == nil {
			return nil, ErrUnknownPkg
		}
		b.Packages[scopeID] = pkg
		cached = pkg
	}

	if err := sn.SemanticAnalysis(cached, &scope.T1SourceMap); err != nil {
		return nil, err
	}
	return fn(source, data)
}

func asNode(v any) (Node, bool) {
	if s, ok := v.(Node); ok {
		return s, true
	}

	if n, ok := v.(ParsedNode); ok {
		return n.AsNode()
	}
	return nil, false
}

func DoParse(parser LexicalParser, source []byte) []Node {
	lexemes := Lex(source)
	var sugars []Node

	offset := 0
	for offset < len(lexemes) {
		slice := lexemes[offset:]

		if parser.Done(slice) {
			if result, success := parser.Result(); success {
				if sugar, ok := asNode(result); ok {
					sugars = append(sugars, sugar)
				}
			}
		}
		offset += parser.Consumed()
		parser.Reset()
	}
	return sugars
}
