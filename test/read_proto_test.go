package test

import (
	"fmt"
	"github.com/emicklei/proto"
	"os"
	"testing"
)

func TestReadProto(t *testing.T) {
	reader, _ := os.Open("D:\\carlos\\hc_enterprise_server\\biz\\dict\\rpc\\dict_model_gen.proto")
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, _ := parser.Parse()

	proto.Walk(definition,
		proto.WithService(handleService),
		proto.WithMessage(handleMessage))
}

func handleService(s *proto.Service) {
	fmt.Println(s.Name)
}

func handleMessage(m *proto.Message) {
	fmt.Println(m.Name)
}
