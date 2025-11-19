package conseal

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
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
	// Filter out generated files
	var userFiles []*ast.File
	for _, f := range pass.Files {
		if !isGeneratedFile(f) {
			userFiles = append(userFiles, f)
		}
	}

	// If all files are generated, skip this package
	if len(userFiles) == 0 {
		return nil, nil
	}

	inspect := inspector.New(userFiles)

	nodeFilter := []ast.Node{
		(*ast.CompositeLit)(nil),
		(*ast.AssignStmt)(nil),
		(*ast.FuncDecl)(nil),
	}

	inspect.WithStack(nodeFilter, func(n ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return true
		}

		filename := pass.Fset.Position(n.Pos()).Filename
		if c.shouldIgnoreFile(filename) {
			return false
		}

		switch node := n.(type) {
		case *ast.CompositeLit:
			c.checkCompositeLit(node, pass, stack)
		case *ast.AssignStmt:
			c.checkAssignStmt(node, pass, stack)
		}
		return true
	})

	return nil, nil
}

func (c *conseal) isTargetPackage(pkgPath string) bool {
	// If no struct-packages are specified, target all packages
	if len(c.config.StructPackages) == 0 {
		return true
	}
	// If struct-packages are specified, only target matching packages
	for _, pattern := range c.config.StructPackages {
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

func isGeneratedFile(file *ast.File) bool {
	for _, commentGroup := range file.Comments {
		for _, comment := range commentGroup.List {
			text := comment.Text
			if strings.Contains(text, "Code generated") && strings.Contains(text, "DO NOT EDIT") {
				return true
			}
		}
	}
	return false
}

func (c *conseal) checkCompositeLit(lit *ast.CompositeLit, pass *analysis.Pass, stack []ast.Node) {
	tv, ok := pass.TypesInfo.Types[lit]
	if !ok {
		return
	}

	typ := tv.Type
	if ptr, ok := typ.(*types.Pointer); ok {
		typ = ptr.Elem()
	}

	named, ok := typ.(*types.Named)
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

	if c.isInAllowedContext(stack, pass, pkgPath) {
		return
	}

	structName := named.Obj().Name()
	pass.Reportf(
		lit.Pos(),
		"direct construction of struct %s is prohibited, use constructor function",
		structName,
	)
}

func (c *conseal) checkAssignStmt(stmt *ast.AssignStmt, pass *analysis.Pass, stack []ast.Node) {
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

		if c.isInAllowedContext(stack, pass, pkgPath) {
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

func (c *conseal) isInAllowedContext(stack []ast.Node, pass *analysis.Pass, targetPkgPath string) bool {
	if c.config.AllowSamePackage && pass.Pkg.Path() == targetPkgPath {
		return true
	}

	enclosingFunc := c.getEnclosingFunc(stack)
	if enclosingFunc == nil {
		return false
	}

	funcName := enclosingFunc.Name.Name

	return c.isConstructor(funcName)
}

func (c *conseal) getEnclosingFunc(stack []ast.Node) *ast.FuncDecl {
	for i := len(stack) - 1; i >= 0; i-- {
		if fn, ok := stack[i].(*ast.FuncDecl); ok {
			return fn
		}
	}
	return nil
}
