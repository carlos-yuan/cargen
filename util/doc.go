package util

import (
	"go/ast"
	goparser "go/parser"
	"go/token"
	"os"
	"strings"
)

func GetPackages(dir string) (map[string]*ast.Package, error) {
	pkgMap, err := goparser.ParseDir(
		token.NewFileSet(),
		dir,
		func(info os.FileInfo) bool {
			// skip go-test
			return !strings.Contains(info.Name(), "_test.go")
		},
		goparser.ParseComments, // no comment
	)
	if err != nil {
		return nil, err
	}
	return pkgMap, nil
}
