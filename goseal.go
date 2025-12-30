package goseal

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

type goseal struct {
	config *Config
}

func NewAnalyzer(config *Config) *analysis.Analyzer {
	c := &goseal{
		config: config,
	}

	return &analysis.Analyzer{
		Name: "goseal",
		Doc:  "Checks that structs are only constructed via factory functions",
		Run:  c.run,
	}
}

func (c *goseal) run(pass *analysis.Pass) (any, error) {
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

func (c *goseal) isTargetPackage(pkgPath string) bool {
	// If no target-packages are specified, target all packages
	if len(c.config.TargetPackages) == 0 {
		return true
	}
	// If target-packages are specified, only target matching packages
	for _, pattern := range c.config.TargetPackages {
		if pattern.MatchString(pkgPath) {
			return true
		}
	}
	return false
}

func (c *goseal) isExcludedStruct(structName string) bool {
	for _, pattern := range c.config.ExcludeStructs {
		if pattern.MatchString(structName) {
			return true
		}
	}
	return false
}

func (c *goseal) shouldIgnoreFile(filename string) bool {
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

func (c *goseal) checkCompositeLit(lit *ast.CompositeLit, pass *analysis.Pass, stack []ast.Node) {
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

	structName := named.Obj().Name()

	if c.isExcludedStruct(structName) {
		return
	}

	if !c.isInitAllowedByScope(pass.Pkg.Path(), pkgPath) {
		pass.Reportf(
			lit.Pos(),
			"direct construction of struct %s is prohibited outside allowed scope",
			structName,
		)
		return
	}

	if !c.isInAllowedFactory(stack) {
		pass.Reportf(
			lit.Pos(),
			"direct construction of struct %s is prohibited, use allowed factory function",
			structName,
		)
		return
	}
}

func (c *goseal) checkAssignStmt(stmt *ast.AssignStmt, pass *analysis.Pass, stack []ast.Node) {
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

		structName := named.Obj().Name()

		if c.isExcludedStruct(structName) {
			continue
		}

		if !c.isMutationAllowedByScope(pass.Pkg.Path(), pkgPath, stack) {
			fieldName := selector.Sel.Name
			pass.Reportf(
				stmt.Pos(),
				"direct assignment to field %s of struct %s is prohibited outside allowed scope",
				fieldName,
				structName,
			)
		}
	}
}

func (c *goseal) isInAllowedFactory(stack []ast.Node) bool {
	// If factory-names is empty, allow all function names
	if len(c.config.FactoryNames) == 0 {
		return true
	}

	enclosingFunc := c.getEnclosingFunc(stack)
	if enclosingFunc == nil {
		return false
	}

	funcName := enclosingFunc.Name.Name

	for _, pattern := range c.config.FactoryNames {
		if pattern.MatchString(funcName) {
			return true
		}
	}
	return false
}

func (c *goseal) isInitAllowedByScope(currentPkg, structPkg string) bool {
	switch c.config.InitScope {
	case InitScopeAny:
		return true

	case InitScopeInTargetPackages:
		return c.isTargetPackage(currentPkg)

	case InitScopeSamePackage:
		return currentPkg == structPkg

	default:
		return false
	}
}

func (c *goseal) isMutationAllowedByScope(currentPkg, structPkg string, stack []ast.Node) bool {
	switch c.config.MutationScope {
	case MutationScopeAny:
		return true

	case MutationScopeReceiver:
		return c.isInReceiverMethod(stack)

	case MutationScopeSamePackage:
		return currentPkg == structPkg

	case MutationScopeNever:
		return false

	default:
		return false
	}
}

func (c *goseal) isInReceiverMethod(stack []ast.Node) bool {
	enclosingFunc := c.getEnclosingFunc(stack)
	if enclosingFunc == nil {
		return false
	}

	if enclosingFunc.Recv == nil || len(enclosingFunc.Recv.List) == 0 {
		return false
	}

	return true
}

func (c *goseal) getEnclosingFunc(stack []ast.Node) *ast.FuncDecl {
	for i := len(stack) - 1; i >= 0; i-- {
		if fn, ok := stack[i].(*ast.FuncDecl); ok {
			return fn
		}
	}
	return nil
}
