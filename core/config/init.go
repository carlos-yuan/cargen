package config

import (
	"comm/aes"
	fileUtl "comm/file"
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/dig"
	"gopkg.in/yaml.v3"
)

func init() {
	var confPath string
	var rootCmd = &cobra.Command{
		Use:   "",
		Short: "hc_entertprise_server",
		Run: func(cmd *cobra.Command, args []string) {
			confPath, _ = cmd.Flags().GetString("config")
		},
	}
	rootCmd.Flags().StringP("config", "c", "", "config path")
	err := rootCmd.Execute()
	if err != nil {
		println("loading config " + err.Error())
	}

	secret, err := loading(confPath)
	if err != nil {
		println("loading config " + err.Error())
	}
	err = yaml.Unmarshal(secret, &global)
	if err != nil {
		println("read config " + err.Error())
	}
	if global.Project != "" {
		Container = dig.New()
		err = Container.Provide(func() *Config {
			return &global
		})
	}
	if err != nil {
		println("read config " + err.Error())
	}
	global.PrintProjectInfo()
}

const (
	ConfigFileName        = "config_origin.yaml"
	ConfigEncryptFileName = "config.yaml"
)

func loading(path string) (bt []byte, err error) {
	if path == "" { //未指定配置 二进制文件所在路径
		path, _ = fileUtl.GetCurrentDirectory()
		paths, err := fileUtl.GetFilePath(path, ConfigEncryptFileName)
		if err != nil {
			return nil, err
		}
		if len(paths) == 1 {
			path = paths[0]
			bt, err = os.ReadFile(path)
			if err != nil {
				return bt, err
			}
		} else if len(paths) > 1 {
			return nil, errors.New("more than one config file")
		} else {
			return nil, errors.New("empty config file")
		}
	} else {
		bt, err = os.ReadFile(path)
		if err != nil {
			return
		}
	}
	if len(bt) == 0 {
		err = errors.New("empty config file")
		return
	}

	// reading unencrypted file
	if !strings.HasSuffix(path, "yaml") {
		return
	}

	// reading encrypted file
	var conf ConfigFile
	err = yaml.Unmarshal(bt, &conf)
	if err != nil {
		return
	}
	if conf.Secret == "" {
		conf.SecretConfig = os.Getenv("SECRET_CONFIG")
		conf.Secret = os.Getenv("SECRET")
	}
	// decode file
	iv, err := aes.DecryptCBC5(conf.Secret, BaseKey, BaseKey)
	if err != nil {
		println("read config " + err.Error())
		return
	}
	return aes.DecryptCBC5Bytes(conf.SecretConfig, BaseKey, iv)
}
