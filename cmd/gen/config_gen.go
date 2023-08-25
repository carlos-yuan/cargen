package gen

import (
	"log"
	"os"
	"path/filepath"

	"github.com/carlos-yuan/cargen/core/config"
	"github.com/carlos-yuan/cargen/util/aes"
	"github.com/carlos-yuan/cargen/util/fileUtil"
	"gopkg.in/yaml.v2"
)

const (
	ConfigFileName        = "config_origin.yaml"
	ConfigEncryptFileName = "config.yaml"
)

func ConfigGen(genPath string) {
	err := filepath.Walk(genPath, func(path string, info os.FileInfo, err error) error {
		if info.Name() == ConfigFileName {
			b, err := fileUtil.ReadAll(path)
			if err != nil {
				log.Fatal(err)
			}
			secretConf := config.ConfigFile{}
			err = yaml.Unmarshal(b, &secretConf)
			if err != nil {
				panic(err)
			}
			secretConf.SecretConfig, err = aes.EncryptCBC5(b, config.BaseKey, secretConf.Secret)
			if err != nil {
				panic(err)
			}
			secretConf.Secret, err = aes.EncryptCBC5([]byte(secretConf.Secret), config.BaseKey, config.BaseKey)
			if err != nil {
				panic(err)
			}
			b, err = yaml.Marshal(&secretConf)
			if err != nil {
				panic(err)
			}
			out, err := fileUtil.CutPathLast(path, 1)
			if err != nil {
				panic(err)
			}
			out = out + string(os.PathSeparator) + ConfigEncryptFileName
			err = fileUtil.WriteStringFile(out, string(b))
			if err != nil {
				panic(err)
			}
			println("generate config " + out)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
