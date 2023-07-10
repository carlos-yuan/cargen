package openapi

import (
	"bytes"
	"github.com/carlos-yuan/cargen/util"
	"go/ast"
	"go/format"
	"go/token"
	"log"
	"net/http"
	"strings"
)

type Api struct {
	Name         string  `json:"name"`         //接口地址
	Summary      string  `json:"summary"`      //名称
	Description  string  `json:"description"`  //描述
	RequestPath  string  `json:"request_path"` //自定义路径
	Point        string  `json:"point"`        //接口结构体对象名称
	Group        string  `json:"group"`        //接口结构体名称
	HttpMethod   string  `json:"method"`       //接口http方法
	Annotate     string  `json:"annotate"`     //注释
	Path         string  `json:"path"`         //文件地址
	Auth         string  `json:"auth"`         //授权方式
	ResponseType string  `json:"responseType"` //返回类型
	Params       *Struct `json:"params"`       //参数 string为路径 Parameter为对象
	Response     *Struct `json:"result"`       //结果 string为路径 Response为对象
	sct          *Struct
}

func (a *Api) GetOperationId() string {
	return a.sct.pkg.Path + "." + a.Name + "." + a.HttpMethod
}

func (a *Api) GetApiPath() string {
	return a.Group + "." + a.Name
}

func (a *Api) NewStruct() *Struct {
	s := NewStruct(a.sct.pkg)
	return &s
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

// CreateApiParameter 构建API基础参数
func (a *Api) CreateApiParameter(body *ast.BlockStmt) {
	for _, line := range body.List {
		//寻找t.Bind(&params)
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
			//寻找 return t.Success(rsp)
			returns := FindReturnInBlockStmt(line)
			for _, rtn := range returns {
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

func FindReturnInBlockStmt(stmt ast.Stmt) []*ast.ReturnStmt {
	var ret []*ast.ReturnStmt
	switch st := stmt.(type) {
	case *ast.ReturnStmt:
		ret = append(ret, st)
	case *ast.IfStmt:
		if st.Body != nil {
			for _, s := range st.Body.List {
				ret = append(ret, FindReturnInBlockStmt(s)...)
			}
		}
	case *ast.SwitchStmt:
		if st.Body != nil {
			for _, s := range st.Body.List {
				ret = append(ret, FindReturnInBlockStmt(s)...)
			}
		}
	case *ast.ForStmt:
		if st.Body != nil {
			for _, s := range st.Body.List {
				ret = append(ret, FindReturnInBlockStmt(s)...)
			}
		}
	case *ast.TypeSwitchStmt:
		if st.Body != nil {
			for _, s := range st.Body.List {
				ret = append(ret, FindReturnInBlockStmt(s)...)
			}
		}
	case *ast.SelectStmt:
		if st.Body != nil {
			for _, s := range st.Body.List {
				ret = append(ret, FindReturnInBlockStmt(s)...)
			}
		}
	case *ast.RangeStmt:
		if st.Body != nil {
			for _, s := range st.Body.List {
				ret = append(ret, FindReturnInBlockStmt(s)...)
			}
		}
	}
	return ret
}

func (a *Api) GetParameterStruct(expr ast.Expr) *Struct {
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
	if structType != nil {
		s := a.NewStruct()
		s.GetStructFromAstStructType(structType)
		s.Name = structName
		return s
	}
	return nil
}

func (a *Api) GetResponseStruct(expr ast.Expr) {
	switch exp := expr.(type) {
	case *ast.Ident: // t.Success(rsp)
		a.Response = a.NewStruct()
		a.Response.Name = exp.Name
		if exp.Obj != nil {
			if decl, ok := exp.Obj.Decl.(*ast.AssignStmt); ok {
				index := FindRspIndex(exp.Name, decl.Lhs)
				a.Response.MethodMap = &MethodMap{Idx: index, Paths: a.FindRspMethodMap(decl.Rhs)}
			}
		}
	case *ast.SelectorExpr: // t.Success(rsp.List)
		a.Response = a.NewStruct()
		id := exp.X.(*ast.Ident)
		a.Response.Name = id.Name
		if id.Obj != nil {
			if decl, ok := id.Obj.Decl.(*ast.AssignStmt); ok {
				index := FindRspIndex(id.Name, decl.Lhs)
				a.Response.MethodMap = &MethodMap{Idx: index, Paths: a.FindRspMethodMap(decl.Rhs)}
			}
		}
		a.Response.Field = exp.Sel.Name
	default:
		println("other " + printAst(exp))
	}
	if a.Response == nil || a.Response.MethodMap == nil {
		s := a.getParameterStruct(expr)
		if s != nil {
			a.Response = s
		}
	}
}

func (a *Api) getParameterStruct(expr ast.Expr) *Struct {
	var structType *ast.StructType
	var structName string
	switch exp := expr.(type) {
	case *ast.UnaryExpr:
		id := exp.X.(*ast.Ident)
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
	case *ast.Ident:
		if spec, ok := exp.Obj.Decl.(*ast.ValueSpec); ok {
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
		} else if stmt, ok := exp.Obj.Decl.(*ast.AssignStmt); ok {
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
	case *ast.CompositeLit:
		if id, ok := exp.Type.(*ast.Ident); ok {
			if id.Obj != nil {
				structType = id.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType)
			}
			structName = id.Name
		}
	case *ast.CallExpr:
		field := GetExprInfo(exp.Fun)
		if field.Pkg == "" && field.Type == "string" { //string方法
			s := a.NewStruct()
			s.Name = field.Type
			return s
		}
	default:
		println("bind parameter error:" + printAst(expr))
	}
	if structType != nil {
		s := a.NewStruct()
		s.GetStructFromAstStructType(structType)
		s.Name = structName
		return s
	}
	return nil
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
			paths = append(paths, a.sct.pkg.Name)
		}
		if a.sct.pkg.Name == paths[len(paths)-1] {
			paths = append(paths, a.sct.pkg.Path)
		} else {
			for name, path := range a.sct.Imports {
				if name == paths[len(paths)-1] {
					paths = append(paths, path)
				}
			}
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

func GetFieldInfo(field *ast.Field) Field {
	f := GetExprInfo(field.Type)
	if field.Names != nil && len(field.Names) == 1 {
		f.Name = field.Names[0].Name
	}
	return f
}

const (
	ExprStruct = "struct"
	ExprFunc   = "func"
)

func GetExprInfo(expr ast.Expr) Field {
	switch exp := expr.(type) {
	case *ast.Ident:
		return Field{Type: exp.Name}
	case *ast.SelectorExpr:
		return Field{Type: exp.Sel.Name, Pkg: exp.X.(*ast.Ident).Name}
	case *ast.StarExpr:
		return GetExprInfo(exp.X)
	case *ast.ArrayType:
		f := GetExprInfo(exp.Elt)
		f.Array = true
		return f
	case *ast.SliceExpr:
		return GetExprInfo(exp.X)
	case *ast.UnaryExpr:
		return GetExprInfo(exp.X)
	case *ast.StructType:
		return Field{Type: ExprStruct}
	case *ast.FuncType:
		return Field{Type: ExprFunc}
	}
	return Field{}
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

// FillRequestParams 填充请求参数
func (a *Api) FillRequestParams(method *Method) {
	if a.Params == nil {
		return
	}
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
}

// FillResponse 填充返回参数
func (a *Api) FillResponse(method *Method) {
	if a.Response == nil {
		return
	}
	// 寻找字段定义
	if a.Response.Fields == nil && a.Response.Name != "" {

		a.sct.pkg.pkgs.FindInMethodMapParams(a.Response)
	}
	if a.Response != nil {
		method.Responses = make(map[string]Response)
		if baseTypes.CheckIn(a.Response.Name) { //非结构体时
			method.Responses["200"] = Response{Description: a.Response.Des, Content: map[string]Content{"application/json": {Schema: Property{Type: (&Field{Type: a.Response.Name}).GetOpenApiType()}}}}
		} else {
			var fields []Field
			isArray := false
			if a.Response.Field != "" {
				field := a.Response.GetField()
				if field != nil {
					fields = field.Struct.Fields
					isArray = field.Array
				}
			} else {
				if len(a.Response.Fields) > 0 {
					fields = a.Response.Fields
				} else {
					log.Printf("response struct not found: %v", a.Response)
				}
			}
			properties := make(map[string]Property)
			for _, field := range fields {
				for _, property := range field.ToProperty() {
					if field.Name != "" {
						if util.FistIsLower(field.Name) {
							continue
						}
					}
					properties[property.Name] = property
				}
			}
			if len(properties) > 0 {
				schemaName := OpenApiSchemasPrefix + a.Response.Name
				method.api.Components.Schemas[a.Response.Name] = Property{
					Type:       OpenApiTypeObject,
					Properties: properties,
				}
				p := Property{}
				if isArray {
					p.Type = OpenApiTypeArray
					p.Items = &Property{Ref: schemaName}
				} else {
					p.Ref = schemaName
				}
				method.Responses["200"] = Response{Description: a.Response.Name, Content: map[string]Content{"application/json": {Schema: p}}}
			}
		}

	}
}
