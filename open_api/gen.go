package openapi

import (
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
	api := base + "/api"
	//todo 使用gomod的名字 做分组
	files, err := GetFilePath(base, "go.mod")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		println(file)
	}
	goModFilePathData, _ := os.ReadFile(base + "/biz/dict/go.mod")
	modFile, _ := modfile.Parse("go.mod", goModFilePathData, nil)
	println(modFile.Module.Mod.Path)
	fs, err := os.ReadDir(base)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range fs {
		println(f.Name())
	}
	//获取api模块目录
	fs, err = os.ReadDir(api)
	if err != nil {
		log.Fatal(err)
	}
	var pkgs []Package
	for _, f := range fs {
		if f.IsDir() {
			ctl := api + "/" + f.Name() + "/controller"
			fs, err := os.ReadDir(ctl)
			if err != nil {
				log.Fatal(err)
			}
			for _, f := range fs {
				if f.IsDir() {
					packages := GenApi(base, ctl+"/"+f.Name(), CreatePackageOpt{NeedApi: true})
					pkgs = append(pkgs, packages...)
				}
			}
		}
	}
	//查找所有导入的包，加载参数依赖
	var refPkg = make(map[string][]Package)
	var refPkgPath = make(map[string]string)
	for _, pkg := range pkgs {
		for pkg, path := range pkg.Imports {
			refPkgPath[pkg] = path
		}
	}
	for _, path := range refPkgPath {
		refPkg[path] = GenApi(base, base+"/"+path)
	}
	println(len(pkgs))
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

type Field struct {
	Name    string `json:"name"`    //名称
	Type    string `json:"type"`    //类型
	Tag     string `json:"tag"`     //tag
	Pkg     string `json:"pkg"`     //包名 类型为结构体时
	PkgPath string `json:"pkgPath"` //包路径 类型为结构体时
	Comment string `json:"comment"` //注释
	Struct  Struct `json:"struct"`  //是结构体时
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

type Document struct {
	Pkg     []Package
	Imports map[string]Package
}

func GenApi(base, path string, opt ...CreatePackageOpt) []Package {
	importPath := path[len(base)+1:]
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
		list = append(list, CreatePackage(pkg, importPath, opt...))
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

type Package struct {
	Name    string            //包名
	Path    string            //包路径
	Mod     string            //模型名 分组用
	Imports map[string]string //包下所有导入的包信息
	Structs []Struct          //所有结构体信息
	Api     []Api             //所有API接口信息
}

type CreatePackageOpt struct {
	NeedApi    bool
	NeedMethod bool
}

func CreatePackage(p *ast.Package, importPath string, opt ...CreatePackageOpt) Package {
	pkg := Package{Name: p.Name, Path: importPath, Imports: make(map[string]string)}
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
						if len(fc.Recv.List) == 1 {
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

func (pkg *Package) GetStructFromAstStructType(s *ast.StructType) Struct {
	sct := Struct{Pkg: pkg.Name, PkgPath: pkg.Path}
	if s != nil {
		for _, fd := range s.Fields.List {
			f := pkg.FieldFromAstField(fd)
			if f.Type == ExprStruct {
				f.Struct = pkg.GetStructFromAstStructType(fd.Type.(*ast.StructType))
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
	}
	f.Comment = FormatComment(fd.Comment)
	return f
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
