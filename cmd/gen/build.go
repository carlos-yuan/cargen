package gen

import (
	"github.com/carlos-yuan/cargen/util/doc"
	"strings"
	"time"

	"github.com/carlos-yuan/cargen/enum"
	openapi "github.com/carlos-yuan/cargen/open_api"
	"github.com/carlos-yuan/cargen/util/fileUtil"
)

type Config struct {
	Gen       string //生成类型
	Path      string //基础项目路径
	Name      string //服务名
	Des       string //描述
	Version   string //版本
	DbDsn     string //数据库Dsn
	Tables    string //数据表 ,隔开
	DbName    string //库别名 避免名称过长或者库名差异
	DictTable string //字典表名
	DictType  string //字典类型字段名
	DictName  string //字典名称字段名
	DictLabel string //字典标签字段名
	DictValue string //字典值字段名
	Out       string //输出目录
}

const (
	GenGrpc   = "grpc"
	GenDB     = "database"
	GenDoc    = "doc"
	GenRouter = "router"
	GenEnum   = "enum"
	GenConfig = "config"
)

func (c Config) Build() {
	start := time.Now().UnixMilli()
	if c.Path != "" {
		c.Path = fileUtil.FixPathSeparator(c.Path)
	}
	if c.Out != "" {
		c.Out = fileUtil.FixPathSeparator(c.Out)
	}
	switch c.Gen {
	case GenGrpc:
		projectPath := c.Path + "/biz/" + c.Name
		if c.Name != "" {
			ModelToProtobuf(c.Path+"/biz", c.Name, "/pb"+c.Name, c.Path+"/orm/"+c.DbName+"/model", "model")
			KitexGen(c.Name, c.Path)
			CarGen(c.Name, c.DbName, projectPath+"/rpc/kitex_gen/pb"+c.Name, "pb"+c.Name, projectPath+"/service/", "service")
		}
	case GenRouter:
		CreateApiRouter(c.Path)
	case GenDB:
		GormGen(c.Path, c.DbDsn, c.DbName, strings.Split(c.Tables, ","))
		if c.DictTable != "" { //检测是否传入字典表
			enum.GenEnum(c.Path+"/orm/"+c.DbName, c.DictTable, c.DictType, c.DictName, c.DictLabel, c.DictValue, c.DbDsn)
		}
	case GenDoc:
		openapi.GenFromPath(c.Name, c.Des, c.Version, c.Path, c.Out)
	case GenEnum:
		enum.GenEnum(c.Path, c.DictTable, c.DictType, c.DictName, c.DictLabel, c.DictValue, c.DbDsn)
	case GenConfig:
		ConfigGen(c.Path)
	}
	if c.Path != "" {
		doc.GoFmt(c.Path)
	}
	println("Generation time:", time.Now().UnixMilli()-start)
}
