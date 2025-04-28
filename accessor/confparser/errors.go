package confparser

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type ErrorConfig struct {
	Code    int    `yaml:"code"`
	Message string `yaml:"msg"`
}

func (e *ErrorConfig) Error() string {
	return fmt.Sprintf("Error Code: %d, Message: %s", e.Code, e.Message)
}

var Errors []ErrorConfig

func InitErrorConfig(filePath string) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Errorf("read config file [%s] error: %v", filePath, err))
	}
	err = yaml.Unmarshal(file, &Errors)
	if err != nil {
		panic(fmt.Errorf("unmarshal config file [%s] error: %v", filePath, err))
	}
}
