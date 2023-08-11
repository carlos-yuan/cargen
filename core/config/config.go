package config

import (
	redisd "comm/redis"
	"fmt"
	"go.uber.org/dig"
	"strconv"
)

var Container *dig.Container

const BaseKey = "oSTP7QSjwvzQAtbI"

var (
	global Config
)

type ConfigFile struct {
	Secret       string `yaml:"secret"`
	SecretConfig string `yaml:"secretConfig"`
}

func (conf *Config) PrintProjectInfo() {
	println("start: ", conf.Project)
}

const (
	EnvDev  = "dev"
	EnvTest = "test"
	EnvPro  = "pro"
)

// Config 配置参数
type Config struct {
	Project  string                   `yaml:"project"`
	Env      string                   `yaml:"env"`
	LogLevel int8                     `yaml:"logLevel"` //全局日志等级
	Gorm     map[string]Gorm          `yaml:"gorm"`
	Redis    map[string]redisd.Config `yaml:"redis"`
	Grpc     map[string]Grpc          `yaml:"grpc"`
	Sms      Sms                      `yaml:"sms"` //短信
	Web      Web                      `yaml:"web"`
}

// Gorm gorm配置参数
type Gorm struct {
	Debug         bool       `yaml:"debug"`
	LogLvl        int        `yaml:"logLvl"`
	SlowThreshold int        `yaml:"slowThreshold"`
	MaxLifetime   int        `yaml:"maxLifeTime"`
	MaxOpenConns  int        `yaml:"maxOpenConns"`
	MaxIdleConns  int        `yaml:"maxIdleConns"`
	Driver        GormDriver `json:"driver"`
	MySQL         MySQL      `yaml:"mysql"`
}

// MySQL mysql配置参数
type MySQL struct {
	Host       string       `yaml:"host"`
	Port       int          `yaml:"port"`
	User       string       `yaml:"user"`
	Password   string       `yaml:"password"`
	DBName     string       `yaml:"dbName"`
	Parameters string       `yaml:"parameters"`
	Variant    MysqlVariant `yaml:"variant"` //mysql变种
}

type GormDriver string

type MysqlVariant string

const (
	Mysql GormDriver = "mysql"
)

// DSN 数据库连接串
func (a MySQL) DSN() string {
	switch a.Variant { //mysql变种判断

	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		a.User, a.Password, a.Host, a.Port, a.DBName, a.Parameters)
}

type Grpc struct {
	Host      string              `yaml:"host"`
	Port      int32               `yaml:"port"`
	WhiteList map[string][]string `yaml:"whiteList"` //白名单
}

func (g Grpc) Address() string {
	return g.Host + ":" + strconv.Itoa(int(g.Port))
}

type Web struct {
	Host   string `yaml:"host"`
	Port   int32  `yaml:"port"`
	Mode   string `yaml:"mode"`
	Prefix string `yaml:"prefix"`
	Token  Token  `yaml:"token"`
}

func (g Web) Address() string {
	return g.Host + ":" + strconv.Itoa(int(g.Port))
}

type Token struct {
	HeaderName string `yaml:"headerName"`
	HeaderType string `yaml:"headerType"`
	CookieName string `yaml:"cookieName"`
	Key        string `yaml:"key"`
	Alg        string `yaml:"alg"`
	Expire     int64  `yaml:"expire"` //过期时间秒
}

type Sms struct {
	RegionId  string `yaml:"regionId"`
	AccessKey string `yaml:"accessKey"`
	Secret    string `yaml:"secret"`
}