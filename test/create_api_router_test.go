package test

import (
	"github.com/carlos-yuan/cargen/cmd/gen"
	"testing"
)

func TestCreateApiRouter(t *testing.T) {
	gen.CreateApiRouter(gen.Config{Path: "D:\\carlos\\hc_enterprise_server"})
}
