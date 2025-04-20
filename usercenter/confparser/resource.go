package confparser

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var ResourceConfig struct {
	Server struct {
		HTTPPort                  int    `yaml:"httpPort"`
		GRPCPort                  int    `yaml:"grpcPort"`
		TokenSecretKey            string `yaml:"tokenSecretKey"`
		TokenExpirationInSecond   int    `yaml:"tokenExpirationInSecond"`
		CaptchaExpirationInSecond int    `yaml:"captchaExpirationInSecond"`
	} `yaml:"server"`
	Logger struct {
		Level               string `yaml:"level"`
		OutPutPath          string `yaml:"outputPath"`
		MaxSingleFileSizeMB int    `yaml:"maxSingleFileSizeMB"`
		MaxBackups          int    `yaml:"maxBackups"`
		MaxStorageAgeInDays int    `yaml:"maxStorageAgeInDays"`
	}
	Redis struct {
		Host               string `yaml:"host"`
		Port               int    `yaml:"port"`
		Password           string `yaml:"pwd"`
		DB                 int    `yaml:"db"`
		PoolSize           int    `yaml:"poolSize"`
		DialTimeoutSecond  int    `yaml:"dialTimeoutSecond"`
		ReadTimeoutSecond  int    `yaml:"readTimeoutSecond"`
		WriteTimeoutSecond int    `yaml:"writeTimeoutSecond"`
		ConnMaxRetries     int    `yaml:"connMaxRetries"`
		TxMaxRetries       int    `yaml:"txMaxRetries"`
		Lock               struct {
			MaxRetries   int `yaml:"maxRetries"`
			RetryDelayMs int `yaml:"retryDelayMs"`
		} `yaml:"lock"`
	}
	MySQL struct {
		Host                    string `yaml:"host"`
		Port                    int    `yaml:"port"`
		Username                string `yaml:"username"`
		Password                string `yaml:"pwd"`
		DBName                  int    `yaml:"dbName"`
		MaxIdleConns            int    `yaml:"maxIdleConns"`
		MaxOpenConns            int    `yaml:"maxOpenConns"`
		ConnMaxLifetimeInSecond int    `yaml:"connMaxLifetimeInSecond"`
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
