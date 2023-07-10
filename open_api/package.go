package openapi

import (
	"github.com/carlos-yuan/cargen/util"
	"go/ast"
	"io/fs"
	"strings"
	"syscall"
)

func (pkgs *Packages) GenApi(pkgPath, path string) []Package {
	astPkg, err := util.GetPackages(path)
	if err != nil {
		fserr, ok := err.(*fs.PathError)
		if ok {
			if fserr.Err == syscall.ENOTDIR {

			} else {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	var list []Package
	for _, pkg := range astPkg {
		list = append(list, pkgs.CreatePackage(pkg, pkgPath))
	}
	//填充所有项目类的结构体依赖
	return list
}

func (pkgs *Packages) CreatePackage(p *ast.Package, pkgPath string) Package {
	pkg := Package{Name: p.Name, Path: pkgPath, Structs: make(map[string]*Struct), pkgs: pkgs}
	//查找结构体定义
	pkg.SetPkgStruct(p)
	//查找API定义
	pkg.SetPkgApi(p)
	return pkg
}

type Package struct {
	Name    string             //包名
	Path    string             //包路径
	Structs map[string]*Struct //所有结构体信息
	pkgs    *Packages
}

type CreatePackageOpt struct {
	NeedApi    bool
	NeedMethod bool
}

// SetPkgStruct 设置包内的结构体
func (pkg *Package) SetPkgStruct(p *ast.Package) {
	for _, file := range p.Files {
		//找到导入的定义
		var imports = make(map[string]string)
		for _, spec := range file.Imports {
			path := strings.ReplaceAll(spec.Path.Value, `"`, "")
			if spec.Name == nil {
				imports[util.LastName(path)] = path
			} else {
				imports[spec.Name.Name] = path
			}
		}
		if file.Scope != nil {
			for _, obj := range file.Scope.Objects {
				s := Struct{Name: obj.Name, Imports: imports, pkg: pkg}
				if obj.Decl != nil {
					ts, ok := obj.Decl.(*ast.TypeSpec)
					if ok {
						s.Des = FormatComment(ts.Comment)
						switch typ := ts.Type.(type) {
						case *ast.StructType:
							s.GetStructFromAstStructType(typ)
						case *ast.InterfaceType:
							s.GetStructFromAstInterfaceType(typ)
						}
					}
				}
				pkg.Structs[s.Name] = &s
			}
		}
	}
}

// FillPkgRelationStruct 设置包内的关联结构体
func (pkgs *Packages) FillPkgRelationStruct() {
	for i, pkg := range *pkgs {
		for is, sct := range pkg.Structs {
			for id, fd := range sct.Fields {
				if !baseTypes.CheckIn(fd.Type) {
					s := pkgs.FindStruct(fd.PkgPath, fd.Pkg, fd.Type)
					if s != nil {
						(*pkgs)[i].Structs[is].Fields[id].Struct = *s
					}
				}
			}
		}
	}
}

func (pkg *Package) SetPkgApi(p *ast.Package) {
	for fp, file := range p.Files {
		//查找API定义
		for _, decl := range file.Decls {
			if fc, ok := decl.(*ast.FuncDecl); ok {
				if fc.Doc != nil {
					api := Api{Name: fc.Name.Name, Path: fp}
					if fc.Recv != nil && len(fc.Recv.List) == 1 {
						f := GetFieldInfo(fc.Recv.List[0])
						api.Point = f.Name
						api.Group = f.Type
					}
					for _, doc := range fc.Doc.List {
						str := strings.TrimSpace(doc.Text)
						str = strings.TrimPrefix(str, `//`)
						str = strings.TrimSpace(str)
						if api.Summary == "" && strings.Index(str, api.Name) == 0 {
							api.Summary = strings.ReplaceAll(str, api.Name, "")
						} else if api.Annotate == "" && strings.Index(str, "@") == 0 {
							api.Annotate = str[1:]
						} else {
							api.Description += str
						}
					}
					api.AnalysisAnnotate()
					if api.HttpMethod != "" {
						api.sct = pkg.Structs[api.Group]
						api.CreateApiParameter(fc.Body)
						api.sct.Api = append(api.sct.Api, api)
					}
					if api.sct != nil { //结构体方法补充
						method := api.sct.GetStructMethodFuncType(fc.Type)
						method.Name = api.Name
						method.PkgPath = pkg.Path
						method.Pkg = api.Group
						api.sct.Methods = append(api.sct.Methods, method)
					}
				}
			}
		}
	}
}

func (sct *Struct) GetStructFromAstStructType(s *ast.StructType) {
	if s != nil {
		for _, fd := range s.Fields.List {
			f := sct.FieldFromAstField(fd)
			if f.Type == "RegionNodeRsp" {
				println(f.Name)
			}
			if !baseTypes.CheckIn(f.Type) {
				field := GetExprInfo(fd.Type)
				if field.Type == ExprStruct {
					f.Struct.Type = field.Pkg
					f.Struct.GetStructFromAstStructType(fd.Type.(*ast.StructType))
				} else {
					f.Array = field.Array
					f.Type = field.Type
					f.Pkg = field.Pkg
					if f.Pkg != "" {
						f.PkgPath = sct.Imports[f.Pkg]
					} else {
						f.Pkg = sct.pkg.Name
						f.PkgPath = sct.pkg.Path
					}
				}
			}
			sct.Fields = append(sct.Fields, f)
		}
	}
}

func (sct *Struct) GetStructFromAstInterfaceType(s *ast.InterfaceType) {
	for _, m := range s.Methods.List {
		sct.Methods = append(sct.Methods, sct.GetInterfaceMethodField(m))
	}
}

func (sct *Struct) FieldFromInterfaceAstField(fd *ast.Field) Field {
	f := GetFieldInfo(fd)
	if f.Pkg != "" {
		f.PkgPath = sct.Imports[f.Pkg]
	}
	if fd.Tag != nil {
		f.Tag = fd.Tag.Value
		f.In, f.ParamName, f.Validate = GetTagInfo(f.Tag)
	}
	f.Comment = FormatComment(fd.Comment)
	return f
}

// GetInterfaceMethodField 寻找接口字段定义
func (sct *Struct) GetInterfaceMethodField(it *ast.Field) StructMethod {
	sm := StructMethod{}
	if len(it.Names) == 1 {
		sm.Name = it.Names[0].Name
	}
	// todo 这点过滤掉了组合接口 基本上不会存在需要返回值时，是在组合接口中
	if ft, ok := it.Type.(*ast.FuncType); ok {
		if ft.Params != nil {
			for _, p := range ft.Params.List {
				sm.Args = append(sm.Args, sct.FieldFromAstField(p))
			}
		}
		if ft.Results != nil {
			for _, r := range ft.Results.List {
				sm.Returns = append(sm.Returns, sct.FieldFromAstField(r))
			}
		}
	}
	return sm
}

// GetStructMethodFuncType 寻找方法参数定义
func (sct *Struct) GetStructMethodFuncType(ft *ast.FuncType) StructMethod {
	sm := StructMethod{}
	if ft.Params != nil {
		for _, p := range ft.Params.List {
			sm.Args = append(sm.Args, sct.FieldFromAstField(p))
		}
	}
	if ft.Results != nil {
		for _, r := range ft.Results.List {
			sm.Returns = append(sm.Returns, sct.FieldFromAstField(r))
		}
	}
	return sm
}

// FieldFromAstField 寻找基础的结构体字段定义
func (sct *Struct) FieldFromAstField(fd *ast.Field) Field {
	f := GetFieldInfo(fd)
	if f.Pkg != "" {
		f.PkgPath = sct.Imports[f.Pkg]
	}
	if fd.Tag != nil {
		f.Tag = fd.Tag.Value
		f.In, f.ParamName, f.Validate = GetTagInfo(f.Tag)
	}
	f.Comment = FormatComment(fd.Comment)
	return f
}
