package openapi

import (
	"github.com/carlos-yuan/cargen/util"
	"go/ast"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Package struct {
	Name    string             //包名
	Path    string             //包路径
	ModPath string             //模块路径
	Structs map[string]*Struct //所有结构体信息
	astPkg  *ast.Package
	pkgs    *Packages
	config  *Config
}

type CreatePackageOpt struct {
	NeedApi    bool
	NeedMethod bool
}

func (pkg *Package) GetAstPkg() *ast.Package {
	return pkg.astPkg
}

const ConfigFileName = "config_origin.yaml"

func (pkg *Package) FindConfig() {
	if pkg.config != nil {
		return
	}
	err := filepath.Walk(pkg.ModPath, func(path string, info os.FileInfo, err error) error {
		if info.Name() == ConfigFileName {
			b, err := util.ReadAll(path)
			if err != nil {
				log.Println(err)
			}
			pkg.config = &Config{}
			err = yaml.Unmarshal(b, pkg.config)
			if err != nil {
				log.Println(err)
			}
		}
		return nil
	})
	if err != nil {
		println(err.Error())
	}
}

// FindPkgStruct 设置包内的结构体
func (pkg *Package) FindPkgStruct() {
	for _, file := range pkg.astPkg.Files {
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
			for _, decl := range file.Decls {
				if gendecl, ok := decl.(*ast.GenDecl); ok {
					for _, spec := range gendecl.Specs {
						ts, ok := spec.(*ast.TypeSpec)
						if ok {
							s := Struct{Name: ts.Name.Name, Imports: imports, pkg: pkg}
							s.Des = FormatComment(gendecl.Doc)
							if s.Des == "" {
								s.Des = s.Name
							}
							switch typ := ts.Type.(type) {
							case *ast.StructType:
								s.GetStructFromAstStructType(typ)
							case *ast.InterfaceType:
								s.GetStructFromAstInterfaceType(typ)
							}
							pkg.Structs[s.Name] = &s
						}
					}
				}
			}
		}
	}
}

func (pkg *Package) FindPkgApi() {
	for fp, file := range pkg.astPkg.Files {
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
							api.Summary = strings.TrimSpace(strings.ReplaceAll(str, api.Name, ""))
						} else if api.Annotate == "" && strings.Index(str, "@") == 0 {
							api.Annotate = str[1:]
						} else {
							api.Description += str
						}
					}
					api.AnalysisAnnotate()
					if api.HttpMethod != "" {
						api.sct = pkg.Structs[api.Group]
						f := GetExprInfo(fc.Type.Results.List[0].Type)
						response := pkg.pkgs.FindStruct(api.sct.Imports[f.Pkg], f.Pkg, f.Type).Copy()
						api.Response = &response
						api.FindApiParameter(fc.Body)
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
