package test

import (
	"testing"

	openapi "github.com/carlos-yuan/cargen/open_api"
	"github.com/carlos-yuan/cargen/util/fileUtil"
)

func TestOpenApi(t *testing.T) {
	path, err := fileUtil.ProjectPath()
	if err != nil {
		t.Error(err)
	}
	openapi.GenFromPath("demo API", "demo的API", "1.0", "D:\\carlos\\cargen_demo", path+"/test/test_api.json")

	//pkgMap, err := goparser.ParseDir(
	//	token.NewFileSet(),
	//	"$GOPATH/src",
	//	func(info os.FileInfo) bool {
	//		// skip go-test
	//		return !strings.Contains(info.Name(), "_test.go")
	//	},
	//	goparser.ParseComments, // no comment
	//)
	//if err != nil {
	//	return nil, err
	//}
}
