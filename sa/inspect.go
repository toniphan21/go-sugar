package sa

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/lex"
)

func InspectPlaceholderFunc(
	pkg *packages.Package,
	smap *sugar.SourceMap,
	node sugar.Node,
	keyword string,
	handler func(file *ast.File, call *ast.CallExpr) bool,
) {
	for _, file := range pkg.Syntax {
		ast.Inspect(file, func(v ast.Node) bool {
			call, ok := v.(*ast.CallExpr)
			if !ok {
				return true
			}

			// is this a __sugar_[keyword]__(...) call?
			ident, ok := call.Fun.(*ast.Ident)
			if !ok || ident.Name != lex.SugarPlaceholderFuncName(keyword) {
				return true
			}

			pos := pkg.Fset.Position(call.Pos())
			originalPosOffset, found := smap.GoToSugarByOffset(pos.Offset)
			if !found || node.Pos().Offset != originalPosOffset {
				return true
			}

			return handler(file, call)
		})
	}
}
