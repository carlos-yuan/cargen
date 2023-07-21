package carpy

import (
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
}

// FindPkgCopyInfo 查找包下所有拷贝信息
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

}
