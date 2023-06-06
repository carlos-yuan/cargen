package main

import (
	"flag"
	"github.com/carlos-yuan/cargen/gen"
)

func main() {
	conf := gen.Config{}
	flag.StringVar(&conf.Name, "n", "", "name")
	flag.StringVar(&conf.Path, "p", "", "path")
	flag.StringVar(&conf.DbDsn, "d", "", "dsn")
	flag.StringVar(&conf.Tables, "t", "", "tables")
	flag.Parse()
	if conf.Name == "" || conf.Path == "" {
		println("服务名及项目目录不能为空")
		return
	}
	conf.Build()
}
