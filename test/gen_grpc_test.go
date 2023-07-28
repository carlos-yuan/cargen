package test

import (
	carkitex "github.com/carlos-yuan/cargen/kitex"
	"testing"
)

func TestGenGRPC(t *testing.T) {
	carkitex.KitexGen("user", "D:\\carlos\\hc_enterprise_server")
}
