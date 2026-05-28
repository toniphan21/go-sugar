package check

import (
	"bytes"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/lex"
)

type sugarImpl struct {
	pos              sugar.Lexeme
	end              sugar.Lexeme
	checkPos         sugar.Lexeme
	checkEnd         sugar.Lexeme
	identifiers      []string
	enclosingReturns []types.Type
	callReturns      []types.Type
}

var _ sugar.Sugar = (*sugarImpl)(nil)

func (s *sugarImpl) Pos() sugar.Lexeme {
	return s.pos
}

func (s *sugarImpl) End() sugar.Lexeme {
	return s.end
}

func (s *sugarImpl) StructuralTransform(source []byte, _ []sugar.Lexeme) []byte {
	// input:  [...]check ...
	// output: [...]__sugar_check__(...)
	// keep: pos-checkPos -> replace "check " by "__sugar_check__(" -> checkEnd:end -> add ")"
	out := bytes.Buffer{}

	out.Write(source[s.pos.Offset:s.checkPos.Offset])
	out.Write([]byte(lex.SugarPlaceholderFuncName(keyword)))
	out.WriteRune('(')
	out.Write(source[s.checkEnd.Offset:s.end.Offset])
	out.WriteRune(')')

	return out.Bytes()
}

func (s *sugarImpl) SemanticTransformer(source []byte, lexemes []sugar.Lexeme) []byte {
	// do simple things first, we update edge cases later
	// output:

	// if no identifiers:
	// if err := [checkEnd:end); err != nil {<nl>
	// <tab>return err<nl>
	// }<nl>

	// if there is identifiers:
	// [identifiers, ] err := [checkEnd:end)\n
	// if err != nil {<nl>
	// <tab>return err<nl>
	// }<nl>
	out := bytes.Buffer{}

	if len(s.identifiers) == 0 {
		out.WriteString("if err := ")
		out.Write(source[s.checkEnd.Offset:s.end.Offset])
		out.WriteString("; err != nil {\n")
		out.WriteRune('\t')
		out.WriteString("return err\n")
		out.WriteRune('}')
	} else {
		idents := make([]string, len(s.identifiers)+1)
		copy(idents, s.identifiers)
		idents[len(s.identifiers)] = "err"

		out.WriteString(strings.Join(idents, ", "))
		out.WriteString(" := ")
		out.Write(source[s.checkEnd.Offset:s.end.Offset])
		out.WriteRune('\n')
		out.WriteString("if err != nil {\n")
		out.WriteRune('\t')
		out.WriteString("return err\n")
		out.WriteRune('}')
	}

	return out.Bytes()
}

func (s *sugarImpl) SemanticAnalysis(pkg *packages.Package, smap *sugar.SourceMap) error {
	for _, file := range pkg.Syntax {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
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
			if !found || s.pos.Offset != originalPosOffset {
				return true
			}

			// reset semantic data
			s.enclosingReturns = nil
			s.callReturns = nil

			// find the enclosing function
			enclosing := findEnclosingFunc(file, call.Pos())
			if enclosing == nil {
				return true
			}
			if enclosing.Type.Results != nil {
				for _, field := range enclosing.Type.Results.List {
					t := pkg.TypesInfo.TypeOf(field.Type)
					s.enclosingReturns = append(s.enclosingReturns, t)
				}
			}

			// the inner expression is the first arg of __sugar_check__
			innerExpr := call.Args[0]
			tv, ok := pkg.TypesInfo.Types[innerExpr]
			if ok {
				if tuple, ok := tv.Type.(*types.Tuple); ok {
					for i := 0; i < tuple.Len(); i++ {
						s.callReturns = append(s.callReturns, tuple.At(i).Type())
					}
				}
			}
			return true
		})
	}
	return nil
}

func findEnclosingFunc(file *ast.File, pos token.Pos) *ast.FuncDecl {
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
