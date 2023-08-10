package main

import (
	"flag"

	"github.com/carlos-yuan/cargen/cmd/gen"
)

func main() {
	conf := gen.Config{}
	flag.StringVar(&conf.Gen, "g", "grpc", "generate type")
	flag.StringVar(&conf.Name, "n", "", "name")
	flag.StringVar(&conf.Path, "p", "", "path")
	flag.StringVar(&conf.DbDsn, "d", "", "dsn")
	flag.StringVar(&conf.Tables, "t", "", "tables")
	flag.StringVar(&conf.DbName, "db", "", "db name")
	flag.StringVar(&conf.Out, "o", "", "output path")
	flag.StringVar(&conf.Des, "des", "", "description")
	flag.StringVar(&conf.Version, "v", "", "版本")
	flag.StringVar(&conf.DictTable, "dictTable", "", "字典表名")
	flag.StringVar(&conf.DictType, "dictType", "type", "字典名称字段名")
	flag.StringVar(&conf.DictName, "dictName", "name", "字典标签字段名")
	flag.StringVar(&conf.DictValue, "dictValue", "value", "字典值字段名")
	flag.Parse()
	if conf.Path == "" {
		panic("项目目录不能为空")
	}
	conf.Build()
}
