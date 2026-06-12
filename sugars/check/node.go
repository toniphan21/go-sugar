package check

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/lex"
)

type node struct {
	pos              sugar.Lexeme
	end              sugar.Lexeme
	checkPos         sugar.Lexeme
	checkEnd         sugar.Lexeme
	identifiers      []string
	enclosingReturns []types.Type
	callReturns      []types.Type
}

var _ sugar.SemanticNode = (*node)(nil)

func (n *node) Pos() sugar.Lexeme {
	return n.pos
}

func (n *node) End() sugar.Lexeme {
	return n.end
}

func (n *node) SemanticAnalysis(pkg *packages.Package, smap *sugar.SourceMap) error {
	for _, file := range pkg.Syntax {
		ast.Inspect(file, func(v ast.Node) bool {
			call, ok := v.(*ast.CallExpr)
			if !ok {
				return true
			}

			// is this a __sugar_check__(...) call?
			ident, ok := call.Fun.(*ast.Ident)
			if !ok || ident.Name != lex.SugarPlaceholderFuncName(keyword) {
				return true
			}

			pos := pkg.Fset.Position(call.Pos())
			originalPosOffset, found := smap.GoToSugarByOffset(pos.Offset)
			if !found || n.pos.Offset != originalPosOffset {
				return true
			}

			// reset semantic data
			n.enclosingReturns = nil
			n.callReturns = nil

			// find the enclosing function
			enclosing := n.findEnclosingFunc(file, call.Pos())
			if enclosing == nil {
				return true
			}
			if enclosing.Type.Results != nil {
				for _, field := range enclosing.Type.Results.List {
					t := pkg.TypesInfo.TypeOf(field.Type)
					n.enclosingReturns = append(n.enclosingReturns, t)
				}
			}

			// the inner expression is the first arg of __sugar_check__
			innerExpr := call.Args[0]
			tv, ok := pkg.TypesInfo.Types[innerExpr]
			if ok {
				if tuple, ok := tv.Type.(*types.Tuple); ok {
					for i := 0; i < tuple.Len(); i++ {
						n.callReturns = append(n.callReturns, tuple.At(i).Type())
					}
				}
			}
			return true
		})
	}
	return nil
}

func (n *node) findEnclosingFunc(file *ast.File, pos token.Pos) *ast.FuncDecl {
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
