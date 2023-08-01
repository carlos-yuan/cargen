package openapi

// Config 配置参数
type Config struct {
	Web struct {
		Prefix string `yaml:"prefix"`
	} `yaml:"web"`
}
