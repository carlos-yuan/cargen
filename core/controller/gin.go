package ctl

import (
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/carlos-yuan/cargen/core/config"
	e "github.com/carlos-yuan/cargen/core/error"
	"github.com/carlos-yuan/cargen/util/timeUtil"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	// ErrEmptyQueryToken can be thrown if authing with URL Query, the query token variable is empty
	ErrEmptyQueryToken = errors.New("query token is empty")

	// ErrEmptyCookieToken can be thrown if authing with a cookie, the token cokie is empty
	ErrEmptyCookieToken = errors.New("cookie token is empty")

	// ErrEmptyAuthHeader can be thrown if authing with a HTTP header, the Auth header needs to be set
	ErrEmptyAuthHeader = errors.New("auth header is empty")

	// ErrInvalidAuthHeader indicates auth header is invalid, could for example have the wrong Realm name
	ErrInvalidAuthHeader = errors.New("auth header is invalid")
)

var validate = validator.New()

func init() {
	// gin 的bind默认设置了参数校验，此处需要设置为空避免多次绑定不同类型造成验证错误
	binding.Validator = nil
}

var _ ControllerContext = &GinControllerContext{}

// GinControllerContext gin 请求体接口实现
type GinControllerContext struct {
	*gin.Context
	conf  *config.Web
	token Token
}

func (c *GinControllerContext) Bind(params any, opts ...BindOption) {
	var bindings = make(map[BindOption]binding.Binding)
	if len(opts) == 0 {
		bindings = constructor.GetBinding(c.Context.Request.Method+c.Context.FullPath(), params)
	} else {
		bindings = constructor.GetBindingByOptions(opts)
	}
	for _, bind := range bindings {
		var err error
		if bind == nil {
			err = c.Context.ShouldBindUri(params)
		} else {
			err = c.Context.ShouldBindWith(params, bind)
		}
		if err != nil && err.Error() == "EOF" {
			err = nil
			continue
		}
		if err != nil {
			panic(e.ParamsDealError.SetErr(err, "参数获取失败"))
		}
	}
	if err := validate.Struct(params); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			if len(errs) > 0 {
				panic(e.ParamsValidatorError.SetErr(errs[0], "参数验证失败"))
			}
		} else {
			panic(e.ParamsValidatorError.SetErr(err, "参数验证失败"))
		}
	}
}

func (c *GinControllerContext) CheckToken(tk Token) {
	if tk == nil {
		panic(e.AuthorizeError.SetErr(errors.New("token instance not set")))
	}
	c.token = tk
	token, err := c.tokenFromHeader(c.token.GetConfig().HeaderName, c.token.GetConfig().HeaderType)
	if err != nil {
		if err != ErrEmptyAuthHeader {
			panic(e.AuthorizeError.SetErr(err))
		}
		if c.token.GetConfig().CookieName != "" {
			token, err = c.tokenFromCookie(c.token.GetConfig().CookieName)
		}
	}
	if err != nil {
		panic(e.AuthorizeError.SetErr(err))
	}
	err = c.token.Verify(token)
	if err != nil {
		panic(e.AuthorizeError.SetErr(err))
	}
	pl := c.token.GetPayLoad()
	if pl.Expire() < timeUtil.Milli() {
		panic(e.AuthorizeTimeOutError)
	}
}

func (c *GinControllerContext) GetToken() Payload {
	return c.token.GetPayLoad()
}

func (c *GinControllerContext) SetHeader(key, val string) {
	c.Header(key, val)
}

func NewGinContext(conf *config.Web) *GinControllerContext {
	return &GinControllerContext{conf: conf}
}

func (c GinControllerContext) SetContext(ctx context.Context) ControllerContext {
	c.Context = ctx.(*gin.Context)
	return &c
}

func (c *GinControllerContext) GetContext() context.Context {
	return c
}

func (c *GinControllerContext) tokenFromHeader(key, typ string) (string, error) {
	authHeader := c.Request.Header.Get(key)
	if authHeader == "" {
		return "", ErrEmptyAuthHeader
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == typ) {
		return "", ErrInvalidAuthHeader
	}
	return parts[1], nil
}

func (c *GinControllerContext) tokenFromQuery(key string) (string, error) {
	token := c.Query(key)
	if token == "" {
		return "", ErrEmptyQueryToken
	}
	return token, nil
}

func (c *GinControllerContext) tokenFromCookie(key string) (string, error) {
	cookie, _ := c.Cookie(key)
	if cookie == "" {
		return "", ErrEmptyCookieToken
	}
	return cookie, nil
}

type BindOption uint8

const (
	_ BindOption = iota
	json
	xml
	yaml
	form
	query
	uri
	file
)

var constructor = &bindConstructor{}

// bindConstructor gin参数绑定构造器
type bindConstructor struct {
	cache map[string]map[BindOption]binding.Binding
	mux   sync.RWMutex
}

func (e *bindConstructor) GetBindingByOptions(opts []BindOption) map[BindOption]binding.Binding {
	var bds map[BindOption]binding.Binding
	for _, opt := range opts {
		bds[opt] = e.GetBindingByOption(opt)
	}
	return bds
}

func (e *bindConstructor) GetBindingByOption(opt BindOption) binding.Binding {
	switch opt {
	case json:
		return &binding.JSON
	case xml:
		return &binding.XML
	case yaml:
		return &binding.YAML
	case form:
		return &binding.Form
	case query:
		return &binding.Query
	case file:
		return ginFormMultipartFileBinding
	case uri:
		return nil
	}
	return nil
}

func (e *bindConstructor) GetBinding(path string, d any) map[BindOption]binding.Binding {
	bs := e.getBinding(path)
	if bs == nil {
		//重新构建
		bs = e.resolve(d)
		e.setBinding(path, bs)
	}
	return bs
}

func (e *bindConstructor) resolve(d any) map[BindOption]binding.Binding {
	bs := make(map[BindOption]binding.Binding, 0)
	qType := reflect.TypeOf(d).Elem()
	var tag reflect.StructTag
	var ok bool
	for i := 0; i < qType.NumField(); i++ {
		tag = qType.Field(i).Tag
		if _, ok = tag.Lookup("json"); ok {
			bs[json] = e.GetBindingByOption(json)
		}
		if _, ok = tag.Lookup("xml"); ok {
			bs[xml] = e.GetBindingByOption(xml)
		}
		if _, ok = tag.Lookup("yaml"); ok {
			bs[yaml] = e.GetBindingByOption(yaml)
		}
		if _, ok = tag.Lookup("form"); ok {
			bs[form] = e.GetBindingByOption(form)
		}
		if _, ok = tag.Lookup("query"); ok {
			bs[query] = e.GetBindingByOption(query)
		}
		if _, ok = tag.Lookup("file"); ok {
			bs[file] = e.GetBindingByOption(file)
		}
		if _, ok = tag.Lookup("uri"); ok {
			bs[uri] = e.GetBindingByOption(uri)
		}
		if t, ok := tag.Lookup("binding"); ok && strings.Index(t, "dive") > -1 {
			qValue := reflect.ValueOf(d)
			for option, b := range e.resolve(qValue.Field(i)) {
				bs[option] = b
			}
			continue
		}
		if t, ok := tag.Lookup("validate"); ok && strings.Index(t, "dive") > -1 {
			qValue := reflect.ValueOf(d)
			for option, b := range e.resolve(qValue.Field(i)) {
				bs[option] = b
			}
		}
	}
	return bs
}

func (e *bindConstructor) getBinding(name string) map[BindOption]binding.Binding {
	return e.cache[name]
}

func (e *bindConstructor) setBinding(name string, bs map[BindOption]binding.Binding) {
	e.mux.Lock()
	defer e.mux.Unlock()
	if e.cache == nil {
		e.cache = make(map[string]map[BindOption]binding.Binding)
	}
	e.cache[name] = bs
}

const defaultMemory = 32 << 20

var ginFormMultipartFileBinding GinFormMultipartFileBinding

type GinFormMultipartFileBinding struct {
}

func (GinFormMultipartFileBinding) Name() string {
	return "multipart/form-data"
}

var typeArrayByte = reflect.TypeOf([]byte{})
var typeArrayArrayByte = reflect.TypeOf([][]byte{})

func (GinFormMultipartFileBinding) Bind(req *http.Request, ptr any) error {
	if req.MultipartForm == nil {
		if err := req.ParseMultipartForm(defaultMemory); err != nil {
			return err
		}
	}
	if err := req.ParseMultipartForm(defaultMemory); err != nil {
		return err
	}
	ptrVal := reflect.ValueOf(ptr)
	if ptrVal.Kind() != reflect.Ptr {
		return errors.New("multipart file binding only works on pointer types")
	}
	for i := 0; i < ptrVal.Elem().NumField(); i++ {
		field := ptrVal.Elem().Type().Field(i)
		tag, ok := field.Tag.Lookup("file")
		if ok {
			files := req.MultipartForm.File[tag]
			if len(files) == 0 {
				continue
			}
			if field.Type == typeArrayByte {
				f, err := files[0].Open()
				defer f.Close()
				if err != nil {
					return err
				}
				b, err := io.ReadAll(f)
				if err != nil {
					return err
				}
				ptrVal.Elem().Field(i).SetBytes(b)
			} else if field.Type == typeArrayArrayByte {
				byts := make([][]byte, len(files))
				for i, file := range files {
					f, err := file.Open()
					defer f.Close()
					if err != nil {
						return err
					}
					b, err := io.ReadAll(f)
					if err != nil {
						return err
					}
					byts[i] = b
				}
				ptrVal.Elem().Field(i).Set(reflect.ValueOf(byts))
			}
		}
	}
	return nil
}

type GinRegister struct {
	Method  string
	Path    string
	Handles []gin.HandlerFunc
}

type GinRegisterList []GinRegister

// 加载路由
func (r GinRegisterList) LoadRoute(g *gin.Engine) {
	for _, handler := range r {
		switch handler.Method {
		case http.MethodGet:
			g.GET(handler.Path, handler.Handles...)
		case http.MethodPost:
			g.POST(handler.Path, handler.Handles...)
		case http.MethodPut:
			g.PUT(handler.Path, handler.Handles...)
		case http.MethodPatch:
			g.PATCH(handler.Path, handler.Handles...)
		case http.MethodDelete:
			g.DELETE(handler.Path, handler.Handles...)
		}
	}
}
