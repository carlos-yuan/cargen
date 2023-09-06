package carpy

import (
	"fmt"
	openapi "github.com/carlos-yuan/cargen/open_api"
	"go/ast"
)

func Gen(base string) {

	cp := Carpy{}
	cp.Gen(base)

}

const (
	PackageName   = "carpy"
	InterfaceName = "Copy"
)

type Carpy struct {
	cpPkg map[string]openapi.Package
}

// Gen 生成
func (c *Carpy) Gen(base string) {
	pkgs := &openapi.Packages{}
	pkgs.InitPackages(base)
	c.FindCopyPkg(pkgs)
	c.FindCopyStruct()
}

// FindCopyPkg 查找包下所有拷贝信息
func (c *Carpy) FindCopyPkg(pkgs *openapi.Packages) {
	c.cpPkg = make(map[string]openapi.Package)
	for _, pkg := range *pkgs {
		for _, file := range pkg.GetAstPkg().Files {
			if file.Scope == nil {
				continue
			}
			for name, obj := range file.Scope.Objects {
				if d, ok := obj.Decl.(*ast.ValueSpec); ok {
					field := openapi.GetExprInfo(d.Type)
					if field.Pkg == PackageName && field.Type == InterfaceName {
						c.cpPkg[name] = pkg
					}
				}
			}
		}
	}
}

// FindCopyStruct 查找所有需要复制的类型
func (c *Carpy) FindCopyStruct() {
	for name, pkg := range c.cpPkg {
		for _, file := range pkg.GetAstPkg().Files {
			ast.Walk(&visitor{name: name}, file)
			for _, decl := range file.Decls {
				//寻找包内方法体定义
				if cpFun, ok := decl.(*ast.FuncDecl); ok {
					for _, body := range cpFun.Body.List {
						if assign, ok := body.(*ast.AssignStmt); ok {
							if len(assign.Rhs) > 0 {
								if call, ok := assign.Rhs[0].(*ast.CallExpr); ok {
									if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
										if ident, ok := selector.X.(*ast.Ident); ok {
											if ident.Name == name && selector.Sel.Name == InterfaceName {
												pkg.FindPkgStruct()
												print(ident.Name)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

type CopyStructInfo struct {
	to   *openapi.Struct
	form *openapi.Struct
}

type visitor struct {
	name string
}

// 寻找所需的
func (v *visitor) Visit(n ast.Node) ast.Visitor {
	if call, ok := n.(*ast.CallExpr); ok {
		if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := selector.X.(*ast.Ident); ok {
				if ident.Name == v.name && selector.Sel.Name == InterfaceName {
					print(ident.Name)
				}
			}
		}
	}
	if n == nil {
		return nil
	}

	fmt.Printf("%T\n", n)
	return v
}
