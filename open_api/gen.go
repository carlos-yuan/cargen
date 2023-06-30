package openapi

import (
	"github.com/carlos-yuan/cargen/util"
	"go/ast"
	"go/doc"
	"log"
	"net/http"
	"os"
	"strings"
)

// GenFromPath 通过目录生成
func GenFromPath(base string) {
	api := base + "/api"
	//获取api模块目录
	fs, err := os.ReadDir(api)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range fs {
		if f.IsDir() {
			ctl := api + "/" + f.Name() + "/controller"
			fs, err := os.ReadDir(ctl)
			if err != nil {
				log.Fatal(err)
			}
			for _, f := range fs {
				if f.IsDir() {
					genApi(ctl + "/" + f.Name())
				}
			}

		}
	}

}

var docs = make(map[string]doc.Package) //map[路径包名]文档 缓存已加载的文档
var linkPath = make(map[string]string)  //map[路径包名]文档 缓存已加载的文档

type Api struct {
	Name         string `json:"name"`         //接口地址
	Summary      string `json:"summary"`      //名称
	Description  string `json:"description"`  //描述
	Point        string `json:"point"`        //接口结构体对象名称
	Group        string `json:"group"`        //接口结构体名称
	HttpMethod   string `json:"method"`       //接口http方法
	Annotate     string `json:"annotate"`     //注释
	Path         string `json:"path"`         //文件地址
	Auth         string `json:"auth"`         //授权方式
	ResponseType string `json:"responseType"` //返回类型
	Params       Struct `json:"params"`       //参数 string为路径 Parameter为对象
	Response     Struct `json:"result"`       //结果 string为路径 Response为对象
}

type Struct struct {
	Name    string  `json:"name"`    //名称
	Pkg     string  `json:"pkg"`     //包名
	PkgPath string  `json:"pkgPath"` //包路径
	Field   string  `json:"field"`   //所需字段 有时可能需要的是结构体中的字段
	Fields  []Field `json:"fields"`  //字段
}

type Field struct {
	Name      string `json:"name"`      //名称
	Type      string `json:"type"`      //类型
	Validator string `json:"validator"` //参数验证
	Pkg       string `json:"pkg"`       //包名 类型为结构体时
	PkgPath   string `json:"pkgPath"`   //包路径 类型为结构体时
}

const (
	AnnotateSplitChar = "|"
	AuthStart         = "Auth:"

	ResponseTypeJSON  = "json"
	ResponseTypeXML   = "xml"
	ResponseTypeBytes = "bytes"

	AuthTypeJWT     = "JWT"
	AuthTypeOAuth2  = "OAuth2"
	AuthTypeSession = "Session"
	AuthTypeCookie  = "Cookie"
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
				a.HttpMethod += s + " "
			}
		}
		if strings.Index(annotate, AuthStart) == 0 {
			a.Auth += strings.ReplaceAll(annotate, AuthStart, "") + " "
		}
		for _, s := range AuthType {
			if s == annotate {
				a.Auth += s + " "
			}
		}
		for _, s := range ResponseType {
			if s == annotate {
				a.ResponseType = s
				break
			}
		}
		if a.ResponseType == "" {
			a.ResponseType = ResponseTypeJSON
		}
	}
}

func genApi(path string) []Api {
	pkgs, err := util.GetPackages(path)
	if err != nil {
		log.Fatal(err)
	}
	var apiList []Api
	for _, pkg := range pkgs {
		for fp, file := range pkg.Files {
			for _, decl := range file.Decls {
				fc, ok := decl.(*ast.FuncDecl)
				if ok {
					if fc.Doc != nil {
						api := Api{Name: fc.Name.Name, Path: fp}
						if len(fc.Recv.List) == 1 {
							api.Point, api.Group = GetRecvInfo(fc.Recv.List[0])
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
							apiList = append(apiList, api)
						}
					}
				}
			}
		}
		docs := doc.New(pkg, path, doc.AllMethods)
		println(docs)
	}
	return apiList
}

func GetRecvInfo(field *ast.Field) (name string, typ string) {
	name = field.Names[0].Name
	switch exp := field.Type.(type) {
	case *ast.StarExpr:
		typ = exp.X.(*ast.Ident).Name
	case *ast.Ident:
		typ = exp.Name
	}
	return
}
