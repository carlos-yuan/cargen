package test

import (
	"testing"

	openapi "github.com/carlos-yuan/cargen/open_api"
	"github.com/carlos-yuan/cargen/util"
)

func TestOpenApi(t *testing.T) {
	path, err := util.ProjectPath()
	if err != nil {
		t.Error(err)
	}
	openapi.GenFromPath("企业端API", "企业端的API", "1.0", "/media/ysgk/DATA/carlos/hc_enterprise_server", path+"/test/test_api.json")
}
