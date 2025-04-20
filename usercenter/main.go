package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"llm_online_inference/usercenter/pb"
	"llm_online_inference/usercenter/server"

	"llm_online_inference/usercenter/confparser"
	"llm_online_inference/usercenter/resource"
	"llm_online_inference/usercenter/router"
	"net"
)

var (
	// 命令行参数
	errorConfigPath    string
	resourceConfigPath string
)

func parseArgs() {
	flag.StringVar(&errorConfigPath, "errorConf", "conf/errors.yml", "path to error config file")
	flag.StringVar(&resourceConfigPath, "resourceConf", "conf/resource.yml", "path to resource config file")
	flag.Parse()
}

func startUpHTTPServer(r *gin.Engine) {
	go func(r *gin.Engine) {
		if err := r.Run(fmt.Sprintf(":%d", confparser.ResourceConfig.Server.HTTPPort)); err != nil {
			panic(fmt.Errorf("failed to start server: %v", err))
		}
	}(r)
}

func startUpGRPCServer(lis net.Listener) {
	s := grpc.NewServer()

	// 注册grpc服务
	pb.RegisterUserCenterServiceServer(s, &server.UserCenterServer{})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		panic(fmt.Errorf("failed to serve: %v", err))
	}
}

func main() {
	parseArgs()

	confparser.InitErrorConfig(errorConfigPath)
	confparser.InitResourceConfig(resourceConfigPath)
	resource.Init()

	// 初始化HTTP服务
	r := gin.Default()
	router.Init(r)

	// 初始化GPRC服务
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", confparser.ResourceConfig.Server.GRPCPort))
	if err != nil {
		panic(fmt.Errorf("failed to listen: %v", err))
	}

	// 启动服务
	startUpHTTPServer(r)
	startUpGRPCServer(lis)
}
