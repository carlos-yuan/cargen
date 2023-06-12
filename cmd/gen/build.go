package gen

import (
	"strings"
	"time"
)

type Config struct {
	Path   string //基础项目路径
	Name   string //服务名
	DbDsn  string //数据库Dsn
	Tables string //数据表 ,隔开
	DbName string //库别名 避免名称过长或者库名差异
}

func (c Config) Build() {
	projectPath := c.Path + "/" + c.Name
	start := time.Now().UnixMilli()
	if c.DbDsn != "" && c.DbName != "" {
		GormGen(c.Path, c.DbDsn, c.DbName, strings.Split(c.Tables, ","))
	}
	ModelToProtobuf(c.Path, c.Name, "pb", c.Path+"/orm/"+c.DbName+"/model", "model")
	KitexGen(c.Name, c.Path)
	CarGen(c.Name, c.DbName, c.Name+"/rpc/kitex_gen/pb", "pb", projectPath+"/service/", "service")
	println("Generation time:", time.Now().UnixMilli()-start)
}
