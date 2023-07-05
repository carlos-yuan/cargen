package openapi

import (
	"encoding/json"
	"github.com/carlos-yuan/cargen/util"
	"go/ast"
	gobuild "go/build"
	"go/doc"
	"golang.org/x/mod/modfile"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"
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
		log.Fatal(err)
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
			apiTags[p.Name] = p.Name
			for _, a := range p.Api {
				name := "/" + util.FistToLower(a.Group) + "/" + util.FistToLower(a.Name)
				if a.RequestPath != "" {
					name = a.RequestPath
				}
				if api.Paths[name] == nil {
					api.Paths[name] = make(map[string]Method)
				}
				method := Method{Tags: []string{p.Name}, OperationId: p.Path + "." + a.Name + "." + a.HttpMethod, Summary: a.Summary}
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
					log.Fatal("generator documentation error for api " + method.OperationId + " to many request types")
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
					method.RequestBody = RequestBody{Content: map[string]Content{"application/json": {Type: PropertyTypeObject, Schema: Property{Type: PropertyTypeObject, Properties: properties}}}}
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

func (pkgs *Packages) Find(pkgPath string) *Package {
	for _, pkg := range *pkgs {
		if pkg.Path == pkgPath {
			return &pkg
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

var docs = make(map[string]doc.Package) //map[路径包名]文档 缓存已加载的文档
var linkPath = make(map[string]string)  //map[路径包名]文档 缓存已加载的文档

type Api struct {
	Name         string `json:"name"`         //接口地址
	Summary      string `json:"summary"`      //名称
	Description  string `json:"description"`  //描述
	RequestPath  string `json:"request_path"` //自定义路径
	Point        string `json:"point"`        //接口结构体对象名称
	Group        string `json:"group"`        //接口结构体名称
	HttpMethod   string `json:"method"`       //接口http方法
	Annotate     string `json:"annotate"`     //注释
	Path         string `json:"path"`         //文件地址
	Auth         string `json:"auth"`         //授权方式
	ResponseType string `json:"responseType"` //返回类型
	Params       Struct `json:"params"`       //参数 string为路径 Parameter为对象
	Response     Struct `json:"result"`       //结果 string为路径 Response为对象
	pkg          *Package
}

type Struct struct {
	Name      string              `json:"name"`    //名称
	Pkg       string              `json:"pkg"`     //包名
	PkgPath   string              `json:"pkgPath"` //包路径
	Field     string              `json:"field"`   //所需字段 有时可能需要的是结构体中的字段
	Fields    []Field             `json:"fields"`  //字段
	MethodMap map[string][]string //map[方法][]链路 方法的返回值类型查找
}

type StructMethod struct {
	Name    string  //名称
	Pkg     string  //包名
	PkgPath string  //包路径
	Args    []Field //参数
	Returns []Field //返回值
}

const (
	AnnotateSplitChar = "|"
	AuthStart         = "auth:"

	ResponseTypeJSON  = "json"
	ResponseTypeXML   = "xml"
	ResponseTypeBytes = "bytes"

	AuthTypeJWT     = "JWT"
	AuthTypeOAuth2  = "OAuth2"
	AuthTypeSession = "Session"
	AuthTypeCookie  = "Cookie"

	ApiParamsBindMethod = "Bind"
)

var APIMethods = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace}

var AuthType = []string{AuthTypeJWT, AuthTypeOAuth2, AuthTypeSession, AuthTypeCookie}

var ResponseType = []string{ResponseTypeJSON, ResponseTypeXML, ResponseTypeBytes}

func (a *Api) AnalysisAnnotate() {
	annotates := strings.Split(a.Annotate, AnnotateSplitChar)
	for _, annotate := range annotates {
		annotate = strings.TrimSpace(annotate)
		for _, s := range APIMethods {
			if s == annotate {
				a.HttpMethod = s
				break
			}
		}
		if a.HttpMethod == annotate {
			a.HttpMethod = strings.ToLower(a.HttpMethod)
			continue
		}
		if strings.Index(annotate, AuthStart) == 0 {
			a.Auth += strings.ReplaceAll(annotate, AuthStart, "") + " "
			continue
		}
		for _, s := range AuthType {
			if s == annotate {
				a.Auth += s + " "
				continue
			}
		}
		for _, s := range ResponseType {
			if s == annotate {
				a.ResponseType = s
				break
			}
		}
		if a.ResponseType == annotate {
			continue
		}
		if a.ResponseType == "" {
			a.ResponseType = ResponseTypeJSON
		}
		a.RequestPath = annotate
	}
}

type Document struct {
	Pkg     []Package
	Imports map[string]Package
}

func GenApi(pkgPath, path string, opt ...CreatePackageOpt) []Package {
	pkgs, err := util.GetPackages(path)
	if err != nil {
		fserr, ok := err.(*fs.PathError)
		if ok {
			if fserr.Err == syscall.ENOTDIR {

			} else {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}
	var list []Package
	for _, pkg := range pkgs {
		list = append(list, CreatePackage(pkg, pkgPath, opt...))
	}
	return list
}

func getImportPkg(pkg, srcDir string) (string, error) {
	p, err := gobuild.Import(pkg, srcDir, gobuild.FindOnly)
	if err != nil {
		return "", err
	}
	return p.Dir, err

}

func (a *Api) CreateApiParameter(body *ast.BlockStmt) {
	for _, line := range body.List {
		expr, ok := line.(*ast.ExprStmt)
		if ok {
			call, ok := expr.X.(*ast.CallExpr)
			if ok {
				method, ok := call.Fun.(*ast.SelectorExpr)
				if ok {
					id, ok := method.X.(*ast.Ident)
					if ok {
						if id.Name == a.Point && method.Sel.Name == ApiParamsBindMethod {
							if expr, ok := call.Args[0].(*ast.UnaryExpr); ok {
								id := expr.X.(*ast.Ident)
								if id.Obj != nil {
									spec := id.Obj.Decl.(*ast.ValueSpec)
									var structType *ast.StructType
									switch st := spec.Type.(type) {
									case *ast.StructType:
										structType = st
									case *ast.Ident:
										structType = st.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType)
									}
									if structType != nil {
										a.Params = a.pkg.GetStructFromAstStructType(structType)
									}
								}
							} else {
								log.Fatal(a.Group + "." + a.Name + " bind params error need a point")
							}
						}
					}
				}
			}
		}
	}
}

func GetFieldInfo(field *ast.Field) (name, pkg, typ string) {
	if field.Names != nil && len(field.Names) == 1 {
		name = field.Names[0].Name
	}
	pkg, typ = GetExprInfo(field.Type)
	return
}

const (
	ExprStruct = "struct"
)

func GetExprInfo(expr ast.Expr) (pkg, typ string) {
	switch exp := expr.(type) {
	case *ast.Ident:
		typ = exp.Name
	case *ast.SelectorExpr:
		pkg = exp.X.(*ast.Ident).Name
		typ = exp.Sel.Name
	case *ast.StarExpr:
		return GetExprInfo(exp.X)
	case *ast.ArrayType:
		pkg, typ = GetExprInfo(exp.Elt)
		typ = "[]" + typ
	case *ast.SliceExpr:
		return GetExprInfo(exp.X)
	case *ast.UnaryExpr:
		return GetExprInfo(exp.X)
	case *ast.StructType:
		return "", ExprStruct
	}
	return
}

func FormatComment(comments *ast.CommentGroup) string {
	str := ""
	if comments == nil {
		return str
	}
	for _, comment := range comments.List {
		str += strings.TrimSpace(strings.TrimPrefix(comment.Text, "//")) + ","
	}
	if str != "" {
		str = str[:len(str)-1]
	}
	return str
}
