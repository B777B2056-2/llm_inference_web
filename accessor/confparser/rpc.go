package confparser

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

var RPCConfig struct {
	Tokenizer struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"tokenizer"`
	ModelServer struct {
		Host string `yaml:"host"`
		Port uint16 `yaml:"port"`
	} `yaml:"model_server"`
}

func InitRPCConfig(filePath string) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Errorf("read config file [%s] error: %v", filePath, err))
	}
	err = yaml.Unmarshal(file, &Errors)
	if err != nil {
		panic(fmt.Errorf("unmarshal config file [%s] error: %v", filePath, err))
	}
}
