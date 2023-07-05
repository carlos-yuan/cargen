package openapi

import (
	"github.com/carlos-yuan/cargen/util"
	"go/ast"
	"strings"
)

type Package struct {
	Name    string            //包名
	Path    string            //包路径
	Imports map[string]string //包下所有导入的包信息
	Structs []Struct          //所有结构体信息
	Api     []Api             //所有API接口信息
}

type CreatePackageOpt struct {
	NeedApi    bool
	NeedMethod bool
}

func CreatePackage(p *ast.Package, pkgPath string, opt ...CreatePackageOpt) Package {
	pkg := Package{Name: p.Name, Path: pkgPath, Imports: make(map[string]string)}
	for fp, file := range p.Files {
		//找到导入的定义
		for _, spec := range file.Imports {
			path := strings.ReplaceAll(spec.Path.Value, `"`, "")
			if spec.Name == nil {
				pkg.Imports[util.LastName(path)] = path
			} else {
				pkg.Imports[spec.Name.Name] = path
			}
		}
		if file.Scope != nil {
			for _, obj := range file.Scope.Objects {
				s := Struct{Name: obj.Name, Pkg: pkg.Name, PkgPath: pkg.Path}
				if obj.Decl != nil {
					ts, ok := obj.Decl.(*ast.TypeSpec)
					if ok {
						st, ok := ts.Type.(*ast.StructType)
						if ok {
							for _, fd := range st.Fields.List {
								s.Fields = append(s.Fields, pkg.FieldFromAstField(fd))
							}
						}
					}
				}
				pkg.Structs = append(pkg.Structs, s)
			}
		}
		//查找API定义
		if len(opt) == 1 && opt[0].NeedApi {
			for _, decl := range file.Decls {
				fc, ok := decl.(*ast.FuncDecl)
				if ok {
					if fc.Doc != nil {
						api := Api{Name: fc.Name.Name, Path: fp, pkg: &pkg}
						if fc.Recv != nil && len(fc.Recv.List) == 1 {
							api.Point, _, api.Group = GetFieldInfo(fc.Recv.List[0])
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
						api.CreateApiParameter(fc.Body)
						api.AnalysisAnnotate()
						if api.HttpMethod != "" {
							pkg.Api = append(pkg.Api, api)
						}
					}
				}
			}
		}
		//查找Method定义
		if len(opt) == 1 && opt[0].NeedMethod {
			//todo 用于加载返回值依赖
		}
	}
	return pkg
}

func (pkg *Package) GetStructFromAstStructType(s *ast.StructType) Struct {
	sct := Struct{Pkg: pkg.Name, PkgPath: pkg.Path}
	if s != nil {
		for _, fd := range s.Fields.List {
			f := pkg.FieldFromAstField(fd)
			if !baseTypes.CheckIn(f.Type) {
				st, ok := fd.Type.(*ast.StructType)
				if ok {
					f.Struct = pkg.GetStructFromAstStructType(st)
				} else {
					sa, ok := fd.Type.(*ast.ArrayType)
					if ok {
						st, ok = sa.Elt.(*ast.StructType)
						if ok {
							f.Struct = pkg.GetStructFromAstStructType(st)
						}
					} else {
						si, ok := fd.Type.(*ast.Ident)
						if ok {
							ts, ok := si.Obj.Decl.(*ast.TypeSpec)
							if ok {
								f.Struct = pkg.GetStructFromAstStructType(ts.Type.(*ast.StructType))
							}
						}
					}
				}
			}
			sct.Fields = append(sct.Fields, f)
		}
	}
	return sct
}

func (pkg *Package) FieldFromAstField(fd *ast.Field) Field {
	f := Field{}
	f.Name, f.Pkg, f.Type = GetFieldInfo(fd)
	if f.Pkg != "" {
		f.PkgPath = pkg.Imports[f.Pkg]
	}
	if fd.Tag != nil {
		f.Tag = fd.Tag.Value
		f.In, f.ParamName, f.Validate = GetTagInfo(f.Tag)
	}
	f.Comment = FormatComment(fd.Comment)
	return f
}
