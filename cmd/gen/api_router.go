package gen

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	openapi "github.com/carlos-yuan/cargen/open_api"
	"github.com/carlos-yuan/cargen/util"
)

// CreateApiRouter 生成api路由
func CreateApiRouter(genPath string) {
	pkgs := openapi.Packages{}
	pkgs.Init(genPath)
	var routers = make(map[string]string) //map[文件路径]代码
	for _, pkg := range pkgs {
		for _, s := range pkg.Structs {
			sort.Slice(s.Api, func(i, j int) bool {
				return strings.Compare(s.Api[i].GetRequestPath(), s.Api[j].GetRequestPath()) == 1
			})
			if len(s.Api) > 0 {
				path := pkg.ModPath + "/router/" + util.ToSnakeCase(s.Name) + ".gen.go"
				importInfo := util.LastName(pkg.Path)
				if importInfo == pkg.Name {
					importInfo = `"` + pkg.Path + `"`
				} else {
					importInfo += ` "` + pkg.Path + `"`
				}
				var apiWriter bytes.Buffer
				for _, api := range s.Api {
					checkToken := ""
					if api.Auth != "" {
						checkToken = "\n\t\t\t\tt.CheckToken()"
					}
					urlPath := api.GetRequestPathNoGroup()
					if api.Params != nil {
						for _, param := range api.Params.Fields {
							if param.In == openapi.OpenApiInPath {
								urlPath = strings.ReplaceAll(urlPath, "{"+param.ParamName+"}", ":"+param.ParamName)
							}
						}
					}
					apiWriter.WriteString(fmt.Sprintf("\n\t\t\tregister{method: \"%s\", path: prefix + \"%s\", handles: []gin.HandlerFunc{func(ctx *gin.Context) {"+
						"\n\t\t\t\tt := t.SetContext(ctx)"+
						"%s"+ //鉴权
						"\n\t\t\t\tctx.JSON(200, t.%s())"+
						"\n\t\t\t}}},",
						strings.ToUpper(api.HttpMethod),
						urlPath,
						checkToken,
						api.Name,
					))
				}
				routers[path] = fmt.Sprintf(apiRouterTemplate, "\t"+importInfo, pkg.Name+"."+s.Name, apiWriter.String())
			}
		}
	}
	for path, src := range routers {
		err := util.WriteByteFile(path, []byte(src))
		if err != nil {
			panic(err)
		}
	}
}

const apiRouterTemplate = `
package router

import (
	"comm/config"
	ctl "comm/controller"
	"comm/convert"
	"github.com/gin-gonic/gin"
	"reflect"
	"strings"
%s
)

func init() {
	err := config.Container.Invoke(func(t *%s, c *config.Config) {
		t.ControllerContext = ctl.NewGinContext(&c.Web)
		typ := reflect.TypeOf(t).Elem()
		prefix := c.Web.Prefix + strings.ToLower(typ.PkgPath()[:strings.Index(typ.PkgPath(), "/")]) + "/" + convert.FistToLower(typ.Name())
		routerList = append(routerList,
%s
		)
	})
	if err != nil {
		panic(err.Error())
	}
}
`
