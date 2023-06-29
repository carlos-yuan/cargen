package test

import (
	openapi "github.com/carlos-yuan/cargen/open_api"
	"testing"
)

func TestOpenApi(t *testing.T) {
	openapi.GenFromPath("D:\\carlos\\hc_enterprise_server")
}
