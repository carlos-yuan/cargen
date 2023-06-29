package openapi

import (
	"github.com/carlos-yuan/cargen/util"
	"log"
)

// GenFromPath 通过目录生成
func GenFromPath(path string) {
	pkgs, err := util.GetPackages(path + "/api")
	if err != nil {
		log.Fatal(err)
	}
	for _, a := range pkgs {
		println(a.Name)
	}
}
