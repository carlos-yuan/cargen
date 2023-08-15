package main

import (
	"os"

	"github.com/carlos-yuan/cargen/core/config"
	"github.com/carlos-yuan/cargen/util/aes"
	fileUtl "github.com/carlos-yuan/cargen/util/file"

	"log"

	"gopkg.in/yaml.v3"
)

func main() {
	var in string
	var out string
	if len(os.Args) == 3 {
		in = os.Args[1]
		out = os.Args[2]
	}
	b, err := fileUtl.ReadAll(in)
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
	err = fileUtl.WriteStringFile(out, string(b))
	if err != nil {
		panic(err)
	}
}
