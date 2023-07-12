package gen

import (
	openapi "github.com/carlos-yuan/cargen/open_api"
)

// CreateApiRouter 生成api路由
func CreateApiRouter(conf Config) {
	pkgs := openapi.Packages{}
	pkgs.Init(conf.Path)
	for _, pkg := range pkgs {
		println(pkg.Name)
	}
}
