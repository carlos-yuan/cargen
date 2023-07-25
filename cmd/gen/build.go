package gen

import (
	openapi "github.com/carlos-yuan/cargen/open_api"
	"github.com/carlos-yuan/cargen/util"
	"strings"
	"time"
)

type Config struct {
	Gen     string //生成类型
	Path    string //基础项目路径
	Name    string //服务名
	Des     string //描述
	Version string //版本
	DbDsn   string //数据库Dsn
	Tables  string //数据表 ,隔开
	DbName  string //库别名 避免名称过长或者库名差异
	Out     string //输出目录
}

const (
	GenGrpc   = "grpc"
	GenDB     = "database"
	GenDoc    = "doc"
	GenRouter = "router"
)

func (c Config) Build() {
	start := time.Now().UnixMilli()
	if c.Gen == GenGrpc {
		projectPath := c.Path + "/biz/" + c.Name
		if c.Name != "" {
			ModelToProtobuf(c.Path+"/biz", c.Name, "/pb"+c.Name, c.Path+"/orm/"+c.DbName+"/model", "model")
			KitexGen(c.Name, c.Path)
			CarGen(c.Name, c.DbName, projectPath+"/rpc/kitex_gen/pb"+c.Name, "pb"+c.Name, projectPath+"/service/", "service")
		}
	} else if c.Gen == GenDB {
		GormGen(c.Path, c.DbDsn, c.DbName, strings.Split(c.Tables, ","))
	} else if c.Gen == GenDoc {
		openapi.GenFromPath(c.Name, c.Des, c.Version, c.Path, c.Out)
	} else if c.Gen == GenRouter {
		CreateApiRouter(c.Path)
	}
	if c.Path != "" {
		util.GoFmt(c.Path)
	}
	println("Generation time:", time.Now().UnixMilli()-start)
}
