package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"llm_online_inference/scheduler/confparser"
	"llm_online_inference/scheduler/resource"
	"llm_online_inference/scheduler/router"
)

var (
	// 命令行参数
	errorConfigPath    string
	resourceConfigPath string
	rpcConfigPath      string
)

func parseArgs() {
	flag.StringVar(&errorConfigPath, "errorConf", "conf/errors.yml", "path to error config file")
	flag.StringVar(&resourceConfigPath, "resourceConf", "conf/resource.yml", "path to resource config file")
	flag.StringVar(&rpcConfigPath, "rpcConf", "conf/rpc.yml", "path to rpc config file")
	flag.Parse()
}

func main() {
	parseArgs()

	confparser.InitErrorConfig(errorConfigPath)
	confparser.InitResourceConfig(resourceConfigPath)
	confparser.InitRPCConfig(rpcConfigPath)
	resource.Init()

	r := gin.Default()

	router.Init(r)

	if err := r.Run(fmt.Sprintf(":%d", confparser.ResourceConfig.Server.HTTPPort)); err != nil {
		panic(fmt.Errorf("failed to start server: %v", err))
	}
}
