package test

import (
	"github.com/carlos-yuan/cargen/cmd/gen"
	"testing"
)

func TestModelToProtobuf(t *testing.T) {
	g := gen.Config{Gen: gen.GenGrpc, Path: "D:\\carlos\\hc_enterprise_server", DbName: "enterprise", Name: "user"}
	g.Build()
}
