package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type conseal struct {
	config *Config
}

func NewAnalyzer(config *Config) *analysis.Analyzer {
	c := &conseal{
		config: config,
	}

	return &analysis.Analyzer{
		Name: "conseal",
		Doc:  "Checks that structs are only constructed via constructor functions",
		Run:  c.run,
	}
}

func (c *conseal) run(pass *analysis.Pass) (any, error) {
	if len(c.config.Packages) == 0 {
		return nil, nil
	}

	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		if c.shouldIgnoreFile(filename) {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.CompositeLit:
				c.checkCompositeLit(n, pass)
			case *ast.AssignStmt:
				c.checkAssignStmt(n, pass)
			}
			return true
		})
	}

	return nil, nil
}

func (c *conseal) isTargetPackage(pkgPath string) bool {
	if len(c.config.Packages) == 0 {
		return false
	}
	for _, pattern := range c.config.Packages {
		if pattern.MatchString(pkgPath) {
			return true
		}
	}
	return false
}

func (c *conseal) isConstructor(funcName string) bool {
	for _, pattern := range c.config.Constructors {
		if pattern.MatchString(funcName) {
			return true
		}
	}
	return false
}

func (c *conseal) shouldIgnoreFile(filename string) bool {
	for _, pattern := range c.config.IgnoreFiles {
		if pattern.MatchString(filename) {
			return true
		}
	}
	return false
}

func (c *conseal) checkCompositeLit(lit *ast.CompositeLit, pass *analysis.Pass) {
	tv, ok := pass.TypesInfo.Types[lit]
	if !ok {
		return
	}

	named, ok := tv.Type.(*types.Named)
	if !ok {
		return
	}

	_, ok = named.Underlying().(*types.Struct)
	if !ok {
		return
	}

	pkg := named.Obj().Pkg()
	if pkg == nil {
		return
	}

	pkgPath := pkg.Path()

	if !c.isTargetPackage(pkgPath) {
		return
	}

	if c.isInAllowedContext(lit.Pos(), pass, pkgPath) {
		return
	}

	structName := named.Obj().Name()
	pass.Reportf(
		lit.Pos(),
		"direct construction of struct %s is prohibited, use constructor function",
		structName,
	)
}

func (c *conseal) checkAssignStmt(stmt *ast.AssignStmt, pass *analysis.Pass) {
	for _, lhs := range stmt.Lhs {
		selector, ok := lhs.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		tv, ok := pass.TypesInfo.Types[selector.X]
		if !ok {
			continue
		}

		typ := tv.Type
		if ptr, ok := typ.(*types.Pointer); ok {
			typ = ptr.Elem()
		}

		named, ok := typ.(*types.Named)
		if !ok {
			continue
		}

		_, ok = named.Underlying().(*types.Struct)
		if !ok {
			continue
		}

		pkg := named.Obj().Pkg()
		if pkg == nil {
			continue
		}

		pkgPath := pkg.Path()

		if !c.isTargetPackage(pkgPath) {
			continue
		}

		if c.isInAllowedContext(stmt.Pos(), pass, pkgPath) {
			continue
		}

		structName := named.Obj().Name()
		fieldName := selector.Sel.Name
		pass.Reportf(
			stmt.Pos(),
			"direct assignment to field %s of struct %s is prohibited, use constructor function",
			fieldName,
			structName,
		)
	}
}

func (c *conseal) isInAllowedContext(pos token.Pos, pass *analysis.Pass, targetPkgPath string) bool {
	if c.config.AllowSamePackage && pass.Pkg.Path() == targetPkgPath {
		return true
	}

	enclosingFunc := c.getEnclosingFunc(pos, pass)
	if enclosingFunc == nil {
		return false
	}

	funcName := enclosingFunc.Name.Name

	return c.isConstructor(funcName)
}

func (c *conseal) getEnclosingFunc(pos token.Pos, pass *analysis.Pass) *ast.FuncDecl {
	var enclosingFunc *ast.FuncDecl

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if n == nil {
				return false
			}

			if pos < n.Pos() || pos > n.End() {
				return false
			}

			if fn, ok := n.(*ast.FuncDecl); ok {
				enclosingFunc = fn
				return false
			}

			return true
		})

		if enclosingFunc != nil {
			break
		}
	}

	return enclosingFunc
}
