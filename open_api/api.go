package openapi

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"net/http"
	"strings"
)

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

func (a *Api) GetOperationId() string {
	return a.pkg.Path + "." + a.Name + "." + a.HttpMethod
}

func (a *Api) GetApiPath() string {
	return a.Group + "." + a.Name
}

type Struct struct {
	Name      string           `json:"name"`      //名称
	Pkg       string           `json:"pkg"`       //包名
	PkgPath   string           `json:"pkgPath"`   //包路径
	Field     string           `json:"field"`     //所需字段 有时可能需要的是结构体中的字段 如：rsp.List
	Fields    []Field          `json:"fields"`    //字段
	MethodMap map[int][]string `json:"methodMap"` //map[参数位置][]链路 方法的返回值类型查找
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

	ApiParamsBindMethod    = "Bind"
	ApiParamsSuccessMethod = "Success"
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

func (a *Api) CreateApiParameter(body *ast.BlockStmt) {
	for _, line := range body.List {
		//寻找t.Bind
		expr, ok := line.(*ast.ExprStmt)
		if ok {
			if call, ok := expr.X.(*ast.CallExpr); ok {
				if method, ok := call.Fun.(*ast.SelectorExpr); ok {
					if id, ok := method.X.(*ast.Ident); ok {
						if id.Name == a.Point && method.Sel.Name == ApiParamsBindMethod {
							a.Params = a.GetParameterStruct(call.Args[0])
						}
					}
				}
			}
		} else {
			//寻找 return t.Success
			if rtn, ok := line.(*ast.ReturnStmt); ok {
				if call, ok := rtn.Results[0].(*ast.CallExpr); ok {
					if method, ok := call.Fun.(*ast.SelectorExpr); ok {
						if id, ok := method.X.(*ast.Ident); ok {
							if id.Name == a.Point && method.Sel.Name == ApiParamsSuccessMethod {
								a.GetResponseStruct(call.Args[0])
							}
						}
					}
				}
			}
		}
	}
}

func (a *Api) GetParameterStruct(expr ast.Expr) Struct {
	var structType *ast.StructType
	var structName string
	if unary, ok := expr.(*ast.UnaryExpr); ok { //形式 &params
		id := unary.X.(*ast.Ident)
		if id.Obj != nil {
			if spec, ok := id.Obj.Decl.(*ast.ValueSpec); ok {
				switch st := spec.Type.(type) {
				case *ast.StructType: //var params struct {}
					structType = st
				case *ast.Ident: //var params Params
					structType = st.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType)
					structName = st.Obj.Decl.(*ast.TypeSpec).Name.Name
				}
			} else if stmt, ok := id.Obj.Decl.(*ast.AssignStmt); ok {
				if com, ok := stmt.Rhs[0].(*ast.CompositeLit); ok {
					if id, ok := com.Type.(*ast.Ident); ok {
						if id.Obj != nil {
							structType = id.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType)
						}
						structName = id.Name
					}
				}
			}
		}
	} else if expr, ok := expr.(*ast.Ident); ok { //形式 params
		if spec, ok := expr.Obj.Decl.(*ast.ValueSpec); ok {
			if spec.Type != nil { // var params Params 参数非指针是不被允许的
				panic("parameter error, please check if the parameter is a pointer or nil:" + printAst(spec))
			} else if spec.Values != nil {
				if unary, ok := spec.Values[0].(*ast.UnaryExpr); ok {
					if com, ok := unary.X.(*ast.CompositeLit); ok {
						switch st := com.Type.(type) {
						case *ast.StructType: //var params = &struct {}{}
							structType = st
						case *ast.Ident: //var params = &DictDataByTypesReq{}
							structType = st.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType)
							structName = st.Name
						}
					}
				}
			}
		} else if stmt, ok := expr.Obj.Decl.(*ast.AssignStmt); ok { //params := &Params{}
			if unary, ok := stmt.Rhs[0].(*ast.UnaryExpr); ok {
				if com, ok := unary.X.(*ast.CompositeLit); ok {
					if id, ok := com.Type.(*ast.Ident); ok {
						if id.Obj != nil {
							structType = id.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType)
						}
						structName = id.Name
					}
				}
			}
		}
	} else {
		panic("bind parameter error:" + printAst(expr))
	}
	var s Struct
	if structType != nil {
		s = a.pkg.GetStructFromAstStructType(structType)
		s.Name = structName
	}
	return s
}

func (a *Api) GetResponseStruct(expr ast.Expr) {
	switch expr := expr.(type) {
	case *ast.Ident: // t.Success(rsp)
		a.Response.Name = expr.Name
		if expr.Obj != nil {
			if decl, ok := expr.Obj.Decl.(*ast.AssignStmt); ok {
				index := FindRspIndex(expr.Name, decl.Lhs)
				a.Response.MethodMap = make(map[int][]string)
				a.Response.MethodMap[index] = a.FindRspMethodMap(decl.Rhs)
			}
		}
		a.Response.Field = expr.Name
	case *ast.SelectorExpr: // t.Success(rsp.List)
		println(expr)
		id := expr.X.(*ast.Ident)
		a.Response.Name = id.Name
		if id.Obj != nil {
			if decl, ok := id.Obj.Decl.(*ast.AssignStmt); ok {
				index := FindRspIndex(id.Name, decl.Lhs)
				a.Response.MethodMap = make(map[int][]string)
				a.Response.MethodMap[index] = a.FindRspMethodMap(decl.Rhs)
			}
		}
		a.Response.Field = expr.Sel.Name
	case *ast.UnaryExpr: // t.Success(&rsp)
		println(expr)
	case *ast.StarExpr: // t.Success(*rsp)
		println(expr)
	default:
		println("other " + printAst(expr))
	}
	if a.Response.MethodMap == nil {
		a.Response = a.getParameterStruct(expr)
	}
}

func (a *Api) getParameterStruct(expr ast.Expr) Struct {
	var structType *ast.StructType
	var structName string
	if unary, ok := expr.(*ast.UnaryExpr); ok {
		id := unary.X.(*ast.Ident)
		if id.Obj != nil {
			if spec, ok := id.Obj.Decl.(*ast.ValueSpec); ok {
				switch st := spec.Type.(type) {
				case *ast.StructType:
					structType = st
				case *ast.Ident:
					structType = st.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType)
					structName = st.Obj.Decl.(*ast.TypeSpec).Name.Name
				}
			} else if stmt, ok := id.Obj.Decl.(*ast.AssignStmt); ok {
				if com, ok := stmt.Rhs[0].(*ast.CompositeLit); ok {
					if id, ok := com.Type.(*ast.Ident); ok {
						if id.Obj != nil {
							structType = id.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType)
						}
						structName = id.Name
					}
				}
			}
		}
	} else if expr, ok := expr.(*ast.Ident); ok {
		if spec, ok := expr.Obj.Decl.(*ast.ValueSpec); ok {
			if spec.Type != nil {
				switch st := spec.Type.(type) {
				case *ast.StructType:
					structType = st
				case *ast.Ident:
					structType = st.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType)
					structName = st.Name
				}
			} else if spec.Values != nil {
				if unary, ok := spec.Values[0].(*ast.UnaryExpr); ok {
					if com, ok := unary.X.(*ast.CompositeLit); ok {
						switch st := com.Type.(type) {
						case *ast.StructType:
							structType = st
						case *ast.Ident:
							structType = st.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType)
							structName = st.Name
						}
					}
				}
			}
		} else if stmt, ok := expr.Obj.Decl.(*ast.AssignStmt); ok {
			if unary, ok := stmt.Rhs[0].(*ast.UnaryExpr); ok {
				if com, ok := unary.X.(*ast.CompositeLit); ok {
					if id, ok := com.Type.(*ast.Ident); ok {
						if id.Obj != nil {
							structType = id.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType)
						}
						structName = id.Name
					}
				}
			}
		}
	} else {
		println("bind parameter error:" + printAst(expr))
	}
	var s Struct
	if structType != nil {
		s = a.pkg.GetStructFromAstStructType(structType)
		s.Name = structName
	}
	return s
}

func (a *Api) FindRspMethodMap(rhs []ast.Expr) []string {
	var paths []string
	for _, hs := range rhs {
		switch expr := hs.(type) {
		case *ast.CallExpr:
			switch expr := expr.Fun.(type) {
			case *ast.SelectorExpr:
				paths = FindSelectorPath(expr)
			case *ast.Ident:
				paths = append(paths, expr.Name)
			default:
				panic("other " + printAst(expr))
			}
		default:
			panic("other " + printAst(expr))
		}
	}
	if len(paths) > 0 {
		if a.Point == paths[len(paths)-1] {
			paths[len(paths)-1] = a.Group
		}
	}
	return paths
}

func FindSelectorPath(s *ast.SelectorExpr) []string {
	var paths = []string{s.Sel.Name}
	if se, ok := s.X.(*ast.SelectorExpr); ok {
		if se.X != nil {
			ps := FindSelectorPath(se)
			paths = append(paths, ps...)
		}
	} else {
		if i, ok := s.X.(*ast.Ident); ok {
			paths = append(paths, i.Name)
		}
	}
	return paths
}

func FindRspIndex(rspName string, lhs []ast.Expr) int {
	for i, hs := range lhs {
		switch expr := hs.(type) {
		case *ast.Ident:
			if expr.Name == rspName {
				return i
			}
		case *ast.StarExpr:
			if expr, ok := expr.X.(*ast.Ident); ok {
				if expr.Name == rspName {
					return i
				}
			}
		}
	}
	return -1
}

func printAst(node interface{}) string {
	defer func() {
		r := recover()
		println(r)
	}()
	var dst bytes.Buffer
	addNewline := func() {
		err := dst.WriteByte('\n') // add newline
		if err != nil {
			panic(err)
		}
	}

	addNewline()

	err := format.Node(&dst, token.NewFileSet(), node)
	if err != nil {
		println(err.Error())
	}

	addNewline()

	return dst.String()
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
