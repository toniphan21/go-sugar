package sa

import (
	"go/ast"
	"go/token"
	"go/types"
)

func FindEnclosingFunc(file *ast.File, pos token.Pos) *ast.FuncDecl {
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fn.Pos() <= pos && pos < fn.End() {
			return fn
		}
	}
	return nil
}

// IsTestingTB returns true if given t is testing.TB
func IsTestingTB(t types.Type) bool {
	if named, isNamed := t.(*types.Named); isNamed {
		if obj := named.Obj(); obj.Pkg() != nil &&
			obj.Pkg().Path() == "testing" && obj.Name() == "TB" {
			return true
		}
	}
	return false
}

// IsTestingReceiver returns true if given t is *testing.T, *testing.B or *testing.F
func IsTestingReceiver(t types.Type) bool {
	ptr, isPtr := t.(*types.Pointer)
	if !isPtr {
		return false
	}
	named, isNamed := ptr.Elem().(*types.Named)
	if !isNamed {
		return false
	}
	obj := named.Obj()
	if obj.Pkg() == nil || obj.Pkg().Path() != "testing" {
		return false
	}
	switch obj.Name() {
	case "T", "B", "F", "TB":
		return true
	default:
		return false
	}
}
