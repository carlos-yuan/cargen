package openapi

import (
	"github.com/carlos-yuan/cargen/util"
	"log"
	"os"
	"path/filepath"
)

// GenFromPath 通过目录生成
func GenFromPath(base string) {
	api := base + "/api"
	err := filepath.Walk(api, func(path string, info os.FileInfo, err error) error {
		println(path)
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
	pkgs, err := util.GetPackages(base + "/api")
	if err != nil {
		log.Fatal(err)
	}
	for _, a := range pkgs {
		println(a.Name)
	}
}
