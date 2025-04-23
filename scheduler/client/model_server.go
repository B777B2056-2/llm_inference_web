package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"llm_online_inference/scheduler/confparser"
	"llm_online_inference/scheduler/pb"
)

type ModelServer struct {
	addr string
}

func NewModelServer() *ModelServer {
	return &ModelServer{
		addr: fmt.Sprintf("%s:%d", confparser.RPCConfig.ModelServer.Host, confparser.RPCConfig.ModelServer.Port),
	}
}

func (m *ModelServer) ChatCompletion(ctx context.Context, inputIDs, attentionMasks []uint32, chatSessionID string) (
	grpc.ServerStreamingClient[pb.ChatCompletionResult], error) {
	conn, err := grpc.NewClient(m.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()
	clt := pb.NewModelServerServiceClient(conn)
	stream, err := clt.ChatCompletion(ctx, &pb.ChatCompletionReq{
		InputIds: inputIDs, AttentionMask: attentionMasks, ChatSessionId: chatSessionID,
	})
	if err != nil {
		return nil, err
	}
	return stream, nil
}
