package client

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"llm_online_interence/llmgateway/confparser"
)

func newGRPCConn(conf confparser.GRPCConfigItem) (*grpc.ClientConn, error) {
	addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
