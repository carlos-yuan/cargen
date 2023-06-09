package openapi

import (
	"encoding/json"
	"github.com/carlos-yuan/cargen/util"
	"golang.org/x/mod/modfile"
	"io/fs"
	"log"
	"os"
	"strings"
	"syscall"
)

// GenFromPath 通过目录生成
func GenFromPath(name, des, version, path, out string) {
	pkgs := Packages{}
	pkgs.Init(path)
	apis := pkgs.GetApi()
	apis.Info.Title = name
	apis.Info.Description = des
	apis.Info.Version = version
	b, _ := json.Marshal(apis)
	err := util.WriteByteFile(out, b)
	if err != nil {
		panic(err)
	}
}

type Packages []Package

func (pkgs *Packages) Init(base string) {
	files, err := GetFilePath(base, "go.mod")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		goModFilePathData, _ := os.ReadFile(file)
		modFile, _ := modfile.Parse("go.mod", goModFilePathData, nil)
		path, err := util.CutPathLast(file, 1)
		if err != nil {
			continue
		}
		pathChild, err := GetAllPath(path)
		if err != nil {
			continue
		}
		//获取所有包和结构体定义
		for _, s := range pathChild {
			pkg := modFile.Module.Mod.Path + "/" + s[len(path)+1:]
			pkg = strings.ReplaceAll(pkg, "\\", "/")
			packages := pkgs.GenPackage(pkg, s)
			*pkgs = append(*pkgs, packages...)
		}
	}
	//填充字段为结构体的依赖
	pkgs.FillPkgRelationStruct()
	//查找API定义
	for i := range *pkgs {
		(*pkgs)[i].FindPkgApi()
	}
}

func (pkgs *Packages) GenPackage(pkgPath, path string) []Package {
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
	//查找结构体定义
	for _, pkg := range astPkg {
		p := Package{Name: pkg.Name, Path: pkgPath, Structs: make(map[string]*Struct), astPkg: pkg, pkgs: pkgs}
		p.SetPkgStruct()
		list = append(list, p)
	}
	return list
}

// FillPkgRelationStruct 设置包内的关联结构体
func (pkgs *Packages) FillPkgRelationStruct() {
	for i, pkg := range *pkgs {
		for is, sct := range pkg.Structs {
			for id, fd := range sct.Fields {
				if !baseTypes.CheckIn(fd.Type) {
					(*pkgs)[i].Structs[is].Fields[id].Struct = pkgs.FindStructPtr(fd.PkgPath, fd.Pkg, fd.Type)
				}
			}
		}
	}
}

// GetApi 获取所有API定义
func (pkgs *Packages) GetApi() OpenAPI {
	api := DefaultInfo
	var apiTags = make(map[string]Tag)
	for _, p := range *pkgs {
		for _, s := range p.Structs {
			if len(s.Api) > 0 {
				apiTags[s.Name] = Tag{Name: s.Name, Description: s.Des}
				for _, a := range s.Api {
					name := "/" + util.FistToLower(a.Group) + "/" + util.FistToLower(a.Name)
					if a.RequestPath != "" {
						name = a.RequestPath
					}
					if api.Paths[name] == nil {
						api.Paths[name] = make(map[string]Method)
					}
					method := Method{Tags: []string{s.Name}, OperationId: a.GetOperationId(), Summary: a.Summary, api: &api}
					a.FillRequestParams(&method)
					a.FillResponse(&method)
					api.Paths[name][a.HttpMethod] = method
				}
			}
		}
	}
	for _, tag := range apiTags {
		api.Tags = append(api.Tags, tag)
	}
	return api
}

// FindInMethodMapParams 从调用链获取返回参数
func (pkgs *Packages) FindInMethodMapParams(sct *Struct) {
	if sct == nil {
		return
	}
	m := sct.MethodMap
	if m == nil || len(m.Paths) == 0 { //非方法调用
		return
	}
	for _, p := range *pkgs {
		if p.Path == m.Paths[len(m.Paths)-1] && p.Name == m.Paths[len(m.Paths)-2] { //匹配包名
			if len(m.Paths) > 2 {
				for _, s := range p.Structs {
					if s.Name == m.Paths[len(m.Paths)-3] { //匹配结构体名
						if len(m.Paths) > 3 {
							var sm *StructMethod
							for i, method := range s.Methods {
								if method.Name == m.Paths[len(m.Paths)-4] {
									sm = &s.Methods[i]
								}
							}
							var st *Struct
							if sm == nil { //可能是组合调用 需要继续寻找更深层结构体的方法
								for _, field := range s.Fields {
									st = pkgs.DeepFindParams(m.Idx, m.Paths[:len(m.Paths)-3], field)
									if st != nil {
										break
									}
								}
							} else {
								field := sm.Returns[m.Idx]
								if baseTypes.CheckIn(field.Type) {
									sct.Type = field.Type
								} else {
									st = pkgs.FindStructPtr(field.PkgPath, field.Pkg, field.Name)
								}
							}
							if st != nil {
								sct.Fields = st.Fields
								sct.Des = st.Des
								sct.Name = st.Name
								sct.Type = st.Type
							}
						} else {
							log.Printf("path too short " + m.Paths[len(m.Paths)-1] + "." + m.Paths[len(m.Paths)-2] + "." + m.Paths[len(m.Paths)-3])
						}
						break
					}
				}
			} else {
				log.Printf("path too short " + m.Paths[len(m.Paths)-1] + "." + m.Paths[len(m.Paths)-2])
			}
			break
		}
	}
}

// DeepFindParams 深度查找参数
func (pkgs *Packages) DeepFindParams(idx int, paths []string, field Field) *Struct {
	if baseTypes.CheckIn(field.Type) {
		return nil
	}
	isComplex := field.Name == "" //组合结构体
	s := pkgs.FindStructPtr(field.PkgPath, field.Pkg, field.Type)
	if s.Name == "" {
		return nil
	}
	var m *StructMethod
	for _, method := range s.Methods {
		if isComplex { //组合匹配
			if method.Name == paths[0] {
				m = &method
				break
			}
		} else {
			if field.Name == paths[len(paths)-1] { //字段名匹配
				if method.Name == paths[len(paths)-2] { //匹配字段结构体方法
					m = &method
					break
				}
			}
		}
	}
	if m == nil {
		for _, field := range s.Fields {
			st := pkgs.DeepFindParams(idx, paths, field)
			if st != nil {
				return st
			}
		}
	} else {
		if len(m.Returns) > idx {
			field := m.Returns[idx]
			if baseTypes.CheckIn(field.Type) {
				return &Struct{Type: field.Type}
			} else {
				return pkgs.FindStructPtr(field.PkgPath, field.Pkg, field.Type)
			}
		}
	}
	return nil
}

func (pkgs *Packages) FindStructPtr(path, pkg, name string) *Struct {
	for _, p := range *pkgs {
		if p.Path == path && p.Name == pkg {
			return p.Structs[name]
		}
	}
	return nil
}

func (pkgs *Packages) FindStruct(path, pkg, name string) Struct {
	for _, p := range *pkgs {
		if p.Path == path && p.Name == pkg {
			s := p.Structs[name]
			if s != nil {
				return s.Copy()
			}
		}
	}
	return Struct{}
}

func GetFilePath(filepath string, fileName string) ([]string, error) {
	infos, err := os.ReadDir(filepath)
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(infos))
	for _, info := range infos {
		path := filepath + string(os.PathSeparator) + info.Name()
		if info.IsDir() {
			tmp, err := GetFilePath(path, fileName)
			if err != nil {
				return nil, err
			}
			paths = append(paths, tmp...)
			continue
		}
		if info.Name() == fileName {
			paths = append(paths, path)
		}
	}
	return paths, nil
}

func GetAllPath(path string) ([]string, error) {
	infos, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(infos))
	for _, info := range infos {
		path := path + string(os.PathSeparator) + info.Name()
		if info.IsDir() {
			tmp, err := GetAllPath(path)
			if err != nil {
				return nil, err
			}
			paths = append(paths, tmp...)
			paths = append(paths, path)
			continue
		}
	}
	return paths, nil
}
