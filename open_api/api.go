package openapi

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"net/http"
	"strings"

	"github.com/carlos-yuan/cargen/util"
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
	AuthTo       string  `json:"authTo"`       //授权方式
	ResponseType string  `json:"responseType"` //返回类型
	Params       *Struct `json:"params"`       //参数 string为路径 Parameter为对象
	Response     *Struct `json:"response"`     //返回结构体
	sct          *Struct
}

func (a *Api) GetOperationId() string {
	return strings.ReplaceAll(a.sct.pkg.Path, "/", ".") + "." + a.Name + "." + a.HttpMethod
}

func (a *Api) GetApiPath() string {
	return a.Group + "." + a.Name
}

func (a *Api) GetRequestPath() string {
	prefix := "/"
	if a.sct.pkg.config != nil && len(a.sct.pkg.config.Web.Prefix) > 0 {
		if a.sct.pkg.config.Web.Prefix[0] != '/' {
			prefix += a.sct.pkg.config.Web.Prefix
		} else {
			prefix = a.sct.pkg.config.Web.Prefix
		}
	}
	prefix += a.sct.pkg.Path[:strings.Index(a.sct.pkg.Path, "/")+1]    //包名
	name := util.FistToLower(a.Group) + "/" + util.FistToLower(a.Name) //结构体名+方法名
	if a.RequestPath != "" {
		if a.RequestPath == "-" {
			name = util.FistToLower(a.Group)
		} else {
			name = util.FistToLower(a.Group) + "/" + a.RequestPath
		}
	}
	return prefix + name
}

func (a *Api) GetRequestPathNoGroup() string {
	name := "/" + util.FistToLower(a.Name)
	if a.RequestPath != "" {
		if a.RequestPath == "-" {
			name = ""
		} else {
			name = "/" + a.RequestPath
		}
	}
	return name
}

func (a *Api) NewStruct() *Struct {
	s := NewStruct(a.sct.pkg)
	return &s
}

// SetResponseDataStruct 返回为结构体时
func (a *Api) SetResponseDataStruct(array bool, s *Struct) {
	for i, field := range a.Response.Fields {
		if field.Name == "Data" {
			a.sct.pkg.pkgs.FindInMethodMapParams(s)
			if s != nil {
				if s.Field != "" { //替换掉指定字段
					for _, f := range s.Fields {
						if f.Name == s.Field {
							a.Response.Fields[i] = Field{Name: field.Name, ParamName: field.ParamName, Struct: f.Struct, Array: f.Array, Type: f.Type, Pkg: f.Pkg, PkgPath: f.PkgPath}
							return
						}
					}
				} else {
					a.Response.Fields[i].Struct = s
					a.Response.Fields[i].Type = s.Name
					a.Response.Fields[i].Array = array
				}
			}
			return
		}
	}
}

// SetResponseDataType 设置返回为基础类型时
func (a *Api) SetResponseDataType(array bool, typ string) {
	for i, field := range a.Response.Fields {
		if field.Name == "Data" {
			a.Response.Fields[i].Array = array
			a.Response.Fields[i].Type = typ
			return
		}
	}
}

func (a *Api) GetResponseData() *Struct {
	for i, field := range a.Response.Fields {
		if field.Name == "Data" {
			return a.Response.Fields[i].Struct
		}
	}
	return nil
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
annotate:
	for _, annotate := range annotates {
		annotate = strings.TrimSpace(annotate)
		for _, s := range APIMethods {
			if s == annotate {
				a.HttpMethod = s
				continue annotate
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
			if len(annotate) >= len(s) && s == annotate[:len(s)] {
				a.Auth = s
				as := strings.Split(annotate, ":")
				if len(as) == 2 {
					a.AuthTo = as[1]
				} else {
					a.AuthTo = a.Auth
				}
				continue annotate
			}
		}
		for _, s := range ResponseType {
			if s == annotate {
				a.ResponseType = s
				continue annotate
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

// FindApiParameter 构建API基础参数
func (a *Api) FindApiParameter(body *ast.BlockStmt) {
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
					if st.Obj != nil {
						structType = st.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType)
						structName = st.Obj.Decl.(*ast.TypeSpec).Name.Name
					} else {
						structName = st.Name
					}
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
		s := a.NewStruct()
		if exp.Obj != nil {
			if decl, ok := exp.Obj.Decl.(*ast.AssignStmt); ok {
				index := FindRspIndex(exp.Name, decl.Lhs)
				s.MethodMap = &MethodMap{Idx: index, Paths: a.FindRspMethodMap(decl.Rhs)}
				a.SetResponseDataStruct(false, s)
			} else if vs, ok := exp.Obj.Decl.(*ast.ValueSpec); ok {
				if idt, ok := vs.Type.(*ast.Ident); ok {
					if baseTypes.CheckIn(idt.Name) { //基础类型
						a.SetResponseDataType(false, idt.Name)
					} else { //结构体
						if idt.Obj != nil {
							if ts, ok := idt.Obj.Decl.(*ast.TypeSpec); ok {
								if st, ok := ts.Type.(*ast.StructType); ok {
									s.GetStructFromAstStructType(st)
									s.Name = idt.Name
									a.SetResponseDataStruct(false, s)
								}
							}
						}
					}
				} else if at, ok := vs.Type.(*ast.ArrayType); ok {
					if idt, ok := at.Elt.(*ast.Ident); ok {
						if baseTypes.CheckIn(idt.Name) { //基础类型
							a.SetResponseDataType(true, idt.Name)
						} else { //结构体
							if idt.Obj != nil {
								if ts, ok := idt.Obj.Decl.(*ast.TypeSpec); ok {
									if st, ok := ts.Type.(*ast.StructType); ok {
										s.GetStructFromAstStructType(st)
										s.Name = idt.Name
										a.SetResponseDataStruct(true, s)
									}
								}
							}
						}
					}
				}
			}
		}
	case *ast.SelectorExpr: // t.Success(rsp.List)
		s := a.NewStruct()
		id := exp.X.(*ast.Ident)
		s.Name = id.Name
		if id.Obj != nil {
			if decl, ok := id.Obj.Decl.(*ast.AssignStmt); ok {
				index := FindRspIndex(id.Name, decl.Lhs)
				s.MethodMap = &MethodMap{Idx: index, Paths: a.FindRspMethodMap(decl.Rhs)}
			}
		}
		s.Field = exp.Sel.Name
		a.SetResponseDataStruct(false, s)
	default:

	}
	if a.GetResponseData() == nil || a.GetResponseData().MethodMap == nil {
		s := a.getParameterStruct(expr)
		a.SetResponseDataStruct(false, s)
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
					if st.Obj != nil {
						structType = st.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType)
					}
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
			for _, property := range field.ToProperty(0, 5) {
				properties[property.Name] = property
			}
		}
		for _, field := range combination {
			for _, property := range field.ToProperty(0, 5) {
				properties[property.Name] = property
			}
		}
		prop := Property{Type: PropertyTypeObject, Properties: properties}
		prop.FillRequired()
		method.RequestBody = RequestBody{Content: map[string]Content{"application/json": {Schema: prop}}}
	}
}

// FillResponse 填充返回参数
func (a *Api) FillResponse(method *Method) {
	if a.Response == nil {
		return
	}
	if a.Response != nil {
		method.Responses = make(map[string]Response)
		if baseTypes.CheckIn(a.Response.Name) { //非结构体时
			method.Responses["200"] = Response{Description: a.Response.Des, Content: map[string]Content{"text/plain": {Schema: Property{Type: (&Field{Type: a.Response.Name}).GetOpenApiType()}}}}
		} else {
			pp := a.Response.ToProperty()
			if len(pp.Properties) > 0 {
				schemaName := OpenApiSchemasPrefix
				name := ""
				data := a.GetResponseData()
				if data != nil {
					if !baseTypes.CheckIn(data.Name) {
						name = data.Name
					}
				}
				if name == "" {
					name = a.Group + "." + a.Name + "Rsp"
				}
				schemaName += name
				method.api.Components.Schemas[name] = pp

				p := Property{Ref: schemaName}
				method.Responses["200"] = Response{Description: a.Response.Name, Content: map[string]Content{"application/json": {Schema: p}}}
			}
		}

	}
}

// FillSecurity 填充Security
func (a *Api) FillSecurity(method *Method) {
	if a.Auth == "" {
		return
	}
	if a.Auth == AuthTypeJWT {
		method.Security = []map[string][]string{{a.AuthTo: []string{}}}
	}
}
