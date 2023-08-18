package test

import (
	"testing"

	"github.com/carlos-yuan/cargen/cmd/gen"
)

func TestGenDatabase(t *testing.T) {
	conf := gen.Config{Gen: gen.GenDB, DbName: "shop",
		DbDsn:     `root:123456@tcp(127.0.0.1:3306)/shop?charset=utf8&parseTime=True&loc=Local&timeout=1000ms`,
		Tables:    `user,goods,orders,dict`,
		Path:      "/media/ysgk/DATA/carlos/cargen_demo",
		DictTable: "dict",
		DictType:  "type",
		DictName:  "name",
		DictLabel: "label",
		DictValue: "value",
	}
	conf.Build()
}
