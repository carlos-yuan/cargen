package ctl

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync"

	"github.com/carlos-yuan/cargen/core/config"
	"github.com/carlos-yuan/cargen/core/controller/token"
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
	token token.Token
}

func (c *GinControllerContext) Bind(params any, opts ...BindOption) {
	var bindings = make(map[BindOption]binding.Binding)
	if len(opts) == 0 {
		bindings = constructor.GetBinding(c.Context.Request.Method+c.Context.FullPath(), params)
	} else {
		bindings = constructor.GetBindingByOptions(opts)
	}
	for i, bind := range bindings {
		var err error
		println(i)
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

func (c *GinControllerContext) CheckToken() {
	var err error
	c.token.Token, err = c.jwtFromHeader(c.conf.Token.HeaderName, c.conf.Token.HeaderType)
	if err != nil {
		if err != ErrEmptyAuthHeader {
			panic(e.AuthorizeError.SetErr(err))
		}
		if c.conf.Token.CookieName != "" {
			c.token.Token, err = c.jwtFromCookie(c.conf.Token.CookieName)
		}
	}
	if err != nil {
		panic(e.AuthorizeError.SetErr(err))
	}
	c.token.Type.Alg = c.conf.Token.Alg
	err = c.token.Verify(c.conf.Token.Key)
	if err != nil {
		panic(e.AuthorizeError.SetErr(err))
	}
	pl := c.token.GetPayLoad()
	if pl.Exp < timeUtil.Milli() {
		panic(e.AuthorizeTimeOutError)
	}
	c.Context.Set(token.GrpcTokenStringKey, c.GetToken())
}

func (c *GinControllerContext) GetToken() *token.Payload {
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

func (c *GinControllerContext) jwtFromHeader(key, typ string) (string, error) {
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

func (c *GinControllerContext) jwtFromQuery(key string) (string, error) {
	token := c.Query(key)
	if token == "" {
		return "", ErrEmptyQueryToken
	}
	return token, nil
}

func (c *GinControllerContext) jwtFromCookie(key string) (string, error) {
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
