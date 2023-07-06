package openapi

import (
	"encoding/json"
	"github.com/carlos-yuan/cargen/util"
	"go/doc"
	"golang.org/x/mod/modfile"
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
			packages := GenApi(pkg, s, CreatePackageOpt{NeedApi: true})
			*pkgs = append(*pkgs, packages...)
		}
	}
}

func (pkgs *Packages) FindApi() OpenAPI {
	api := DefaultInfo
	api.Paths = make(map[string]map[string]Method)
	var apiTags = make(map[string]string)
	for _, p := range *pkgs {
		if len(p.Api) > 0 {
			var structApi = make(map[string]string)
			for _, a := range p.Api {
				structApi[a.Group] = a.Group
			}
			for k, _ := range structApi {
				for _, s := range p.Structs {
					if k == s.Name {
						pkgs.GetMethodMap(s)
					}
				}
			}
			apiTags[p.Name] = p.Name
			for _, a := range p.Api {
				name := "/" + util.FistToLower(a.Group) + "/" + util.FistToLower(a.Name)
				if a.RequestPath != "" {
					name = a.RequestPath
				}
				if api.Paths[name] == nil {
					api.Paths[name] = make(map[string]Method)
				}
				method := Method{Tags: []string{p.Name}, OperationId: a.GetOperationId(), Summary: a.Summary}
				var jsonFields []Field
				var xmlFields []Field
				var yamlFields []Field
				var combination []Field //组合参数
				for _, f := range a.Params.Fields {
					switch f.In {
					case TagParamJson:
						jsonFields = append(jsonFields, f)
					case TagParamXml:
						xmlFields = append(xmlFields, f)
					case TagParamYaml:
						yamlFields = append(yamlFields, f)
					default:
						if f.Name == "" {
							combination = append(combination, f)
						} else {
							method.Parameters = append(method.Parameters, f.ToParameter())
						}
					}
				}
				if (len(jsonFields) > 0 && len(xmlFields) > 0) || (len(jsonFields) > 0 && len(yamlFields) > 0) || (len(yamlFields) > 0 && len(xmlFields) > 0) {
					panic("generator documentation error for api " + method.OperationId + " to many request types")
				}
				if len(jsonFields) > 0 {
					properties := make(map[string]Property)
					for _, field := range jsonFields {
						for _, property := range field.ToProperty() {
							properties[property.Name] = property
						}
					}
					for _, field := range combination {
						for _, property := range field.ToProperty() {
							properties[property.Name] = property
						}
					}
					method.RequestBody = RequestBody{Content: map[string]Content{"application/json": {Schema: Property{Type: PropertyTypeObject, Properties: properties}}}}
				}
				api.Paths[name][a.HttpMethod] = method
			}
		}
	}
	for _, tag := range util.MapToSplice(apiTags) {
		api.Tags = append(api.Tags, Tag{Name: tag})
	}
	return api
}

func (pkgs *Packages) GetMethodMap(s Struct) {

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

var docs = make(map[string]doc.Package) //map[路径包名]文档 缓存已加载的文档
var linkPath = make(map[string]string)  //map[路径包名]文档 缓存已加载的文档
