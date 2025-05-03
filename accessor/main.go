package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"llm_inference_web/accessor/confparser"
	"llm_inference_web/accessor/resource"
	"llm_inference_web/accessor/router"
)

var (
	// 命令行参数
	resourceConfigPath string
	rpcConfigPath      string
)

func parseArgs() {
	flag.StringVar(&resourceConfigPath, "resourceConf", "conf/resource.yml", "path to resource config file")
	flag.StringVar(&rpcConfigPath, "rpcConf", "conf/rpc.yml", "path to rpc config file")
	flag.Parse()
}

func main() {
	parseArgs()

	confparser.InitResourceConfig(resourceConfigPath)
	confparser.InitRPCConfig(rpcConfigPath)
	resource.Init()

	r := gin.Default()

	router.Init(r)

	if err := r.Run(fmt.Sprintf(":%d", confparser.ResourceConfig.Server.HTTPPort)); err != nil {
		panic(fmt.Errorf("failed to start server: %v", err))
	}
}
