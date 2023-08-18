package test

import (
	"testing"

	"github.com/carlos-yuan/cargen/cmd/gen"
)

func TestGenGRPC(t *testing.T) {
	conf := gen.Config{Gen: gen.GenGrpc, Path: "/media/ysgk/DATA/carlos/cargen_demo", DbName: "shop", Name: "user"}
	conf.Build()
}
