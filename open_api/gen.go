package openapi

import (
	"encoding/json"
	"github.com/carlos-yuan/cargen/util"
	"golang.org/x/mod/modfile"
	"log"
	"os"
	"strings"
)

// GenFromPath 通过目录生成
func GenFromPath(base string) {
	pkgs := Packages{}
	pkgs.Init(base)
	apis := pkgs.FindApi()
	apis.Info.Title = "hc_enterprise_api"
	apis.Info.Description = "hc_enterprise_api"
	apis.Info.Version = "v0.0.1"
	b, _ := json.Marshal(apis)
	println(string(b))
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
		for _, s := range pathChild {
			pkg := modFile.Module.Mod.Path + "/" + s[len(path)+1:]
			pkg = strings.ReplaceAll(pkg, "\\", "/")
			packages := pkgs.GenApi(pkg, s)
			*pkgs = append(*pkgs, packages...)
		}
	}
}

func (pkgs *Packages) FindApi() OpenAPI {
	api := DefaultInfo
	var apiTags = make(map[string]string)
	for _, p := range *pkgs {
		for _, s := range p.Structs {
			if len(s.Api) > 0 {
				apiTags[p.Name] = p.Name
				for _, a := range s.Api {
					name := "/" + util.FistToLower(a.Group) + "/" + util.FistToLower(a.Name)
					if a.RequestPath != "" {
						name = a.RequestPath
					}
					if api.Paths[name] == nil {
						api.Paths[name] = make(map[string]Method)
					}
					method := Method{Tags: []string{p.Name}, OperationId: a.GetOperationId(), Summary: a.Summary, api: &api}
					a.FillRequestParams(&method)
					a.FillResponse(&method)
					api.Paths[name][a.HttpMethod] = method
				}
			}
		}
	}
	for _, tag := range util.MapToSplice(apiTags) {
		api.Tags = append(api.Tags, Tag{Name: tag})
	}
	return api
}

func (pkgs *Packages) FindParams(m MethodMap) *Struct {
	if len(m.Paths) == 0 {
		return nil
	}
	for _, p := range *pkgs {
		if p.Path == m.Paths[len(m.Paths)-1] && p.Name == m.Paths[len(m.Paths)-2] { //匹配包名
			if len(m.Paths) > 2 {
				for _, s := range p.Structs {
					if s.Name == m.Paths[len(m.Paths)-3] { //匹配结构体名
						if len(m.Paths) > 3 {
							var sm *StructMethod
							for _, method := range s.Methods {
								if method.Name == m.Paths[len(m.Paths)-4] {
									sm = &method
								}
							}
							if sm == nil { //可能是组合调用 需要继续寻找更深层结构体的方法
								for _, field := range s.Fields {
									return pkgs.DeepFindParams(m.Idx, m.Paths[:len(m.Paths)-3], field)
								}
							} else {
								field := sm.Returns[m.Idx]
								if baseTypes.CheckIn(field.Type) {
									return &Struct{Type: field.Type}
								} else {
									return pkgs.FindStruct(field.PkgPath, field.Pkg, field.Name)
								}
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
	return nil
}

func (pkgs *Packages) DeepFindParams(idx int, paths []string, field Field) *Struct {
	if baseTypes.CheckIn(field.Type) {
		return nil
	}
	isComplex := field.Name == "" //组合结构体
	s := pkgs.FindStruct(field.PkgPath, field.Pkg, field.Type)
	if s == nil {
		return nil
	}
	var m *StructMethod
	for _, method := range s.Methods {
		if isComplex { //组合匹配
			if method.Name == paths[0] {
				m = &method
			}
		} else {
			if field.Name == paths[len(paths)-1] { //字段名匹配
				if method.Name == paths[len(paths)-2] { //匹配字段结构体方法
					m = &method
				}
			}
		}
	}
	if m == nil {
		for _, field := range s.Fields {
			return pkgs.DeepFindParams(idx, paths, field)
		}
	} else {
		if len(m.Returns) > idx {
			field := m.Returns[idx]
			if baseTypes.CheckIn(field.Type) {
				return &Struct{Type: field.Type}
			} else {
				return pkgs.FindStruct(field.PkgPath, field.Pkg, field.Type)
			}
		}
	}
	return nil
}

func (pkgs *Packages) FindStruct(path, pkg, name string) *Struct {
	for _, p := range *pkgs {
		if p.Path == path && p.Name == pkg {
			return p.Structs[name]
		}
	}
	return nil
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
