package main

import (
	"flag"
	"fmt"
	"llm_online_interence/llmgateway/confparser"
	"llm_online_interence/llmgateway/limiter"
	"llm_online_interence/llmgateway/middleware"
	"llm_online_interence/llmgateway/proxy"
	"llm_online_interence/llmgateway/resource"

	"github.com/gin-gonic/gin"
)

var (
	// 命令行参数
	proxyConfigPath    string
	resourceConfigPath string
)

func parseArgs() {
	flag.StringVar(&proxyConfigPath, "proxyConf", "conf/proxy.yml", "path to proxy config file")
	flag.StringVar(&resourceConfigPath, "resourceConf", "conf/resource.yml", "path to resource config file")
	flag.Parse()
}

func bootstrap() {
	// 解析启动参数
	parseArgs()
	// 加载配置文件
	confparser.InitProxyConfig(proxyConfigPath)
	confparser.InitResourceConfig(resourceConfigPath)
	// 初始化资源
	resource.Init()
	// 限流器初始化
	limiter.Init()
}

func main() {
	// 初始化服务
	bootstrap()

	r := gin.Default()

	// 注册中间件
	r.Use(
		middleware.TraceID(),    // 链路追踪
		middleware.BlackedIPs(), // 阻止黑名单IP访问
	)

	// 初始化代理服务
	proxy.Init(r)

	// 启动服务
	if err := r.Run(fmt.Sprintf(":%d", confparser.ResourceConfig.Server.Port)); err != nil {
		resource.Logger.Fatal("failed to start server", err)
	}
}
