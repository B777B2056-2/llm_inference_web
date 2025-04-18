package confparser

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var ResourceConfig struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Logger struct {
		Level               string `yaml:"level"`
		OutPutPath          string `yaml:"outputPath"`
		MaxSingleFileSizeMB int    `yaml:"maxSingleFileSizeMB"`
		MaxBackups          int    `yaml:"maxBackups"`
		MaxStorageAgeInDays int    `yaml:"maxStorageAgeInDays"`
	}
	Redis struct {
		Host           string `yaml:"host"`
		Port           int    `yaml:"port"`
		Password       string `yaml:"pwd"`
		DB             int    `yaml:"db"`
		PoolSize       int    `yaml:"poolSize"`
		DialTimeout    int    `yaml:"dialTimeout"`
		ReadTimeout    int    `yaml:"readTimeout"`
		WriteTimeout   int    `yaml:"writeTimeout"`
		ConnMaxRetries int    `yaml:"connMaxRetries"`
		TxMaxRetries   int    `yaml:"txMaxRetries"`
	}
}

func InitResourceConfig(filePath string) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Errorf("read config file [%s] error: %v", filePath, err))
	}
	err = yaml.Unmarshal(file, &ResourceConfig)
	if err != nil {
		panic(fmt.Errorf("unmarshal config file [%s] error: %v", filePath, err))
	}
	ResourceConfig.Logger.Level = strings.ToLower(ResourceConfig.Logger.Level)
}
