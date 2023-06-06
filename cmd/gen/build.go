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
}

func (c Config) Build() {
	projectPath := c.Path + "/" + c.Name
	start := time.Now().UnixMilli()
	if c.DbDsn != "" {
		GormGen(projectPath, c.DbDsn, strings.Split(c.Tables, ","))
	}
	ModelToProtobuf(c.Path, c.Name, "pb", c.Name+"/orm/model", "model")
	KitexGen(c.Name, c.Path)
	CarGen(c.Name, c.Name+"/rpc/kitex_gen/pb", "pb", projectPath+"/service/", "service")
	println("Generation time:", time.Now().UnixMilli()-start)
}
