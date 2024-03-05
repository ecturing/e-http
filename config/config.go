package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var (
	Config *yaml.Decoder
)

func init() {
	// 初始化yaml配置文件器
	file, err := os.Open("config/example_config.yaml")
	if err != nil {
		panic(err)
	}
	Config = yaml.NewDecoder(file)
}
