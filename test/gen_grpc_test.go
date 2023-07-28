package test

import (
	"github.com/carlos-yuan/cargen/cmd/gen"
	"testing"
)

func TestGenGRPC(t *testing.T) {
	gen.KitexGen("user", "D:\\carlos\\hc_enterprise_server")
}
