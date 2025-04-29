package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"llm_online_inference/accessor/confparser"
	"llm_online_inference/accessor/dto"
	"llm_online_inference/accessor/pb"
)

type ModelServer struct {
	addr       string
	streamConn *grpc.ClientConn
}

func NewModelServer() *ModelServer {
	return &ModelServer{
		addr: fmt.Sprintf("%s:%d", confparser.RPCConfig.ModelServer.Host, confparser.RPCConfig.ModelServer.Port),
	}
}

func (m *ModelServer) ChatCompletion(ctx context.Context, chatSessionID string, tokenIds, tokenTypeIds []uint64,
	params dto.InferenceParams) (grpc.ServerStreamingClient[pb.ChatCompletionResult], error) {
	conn, err := grpc.NewClient(m.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	m.streamConn = conn
	clt := pb.NewModelServerServiceClient(conn)
	traceId, ok := ctx.Value("trace_id").(string)
	if !ok {
		traceId = "unknown"
	}
	stream, err := clt.ChatCompletion(ctx, &pb.ChatCompletionReq{
		ChatSessionId:     chatSessionID,
		TokenIds:          tokenIds,
		TokenTypeIds:      tokenTypeIds,
		PresencePenalty:   params.PresencePenalty,
		FrequencyPenalty:  params.FrequencyPenalty,
		RepetitionPenalty: params.RepetitionPenalty,
		Temperature:       params.Temperature,
		TopP:              params.TopP,
		TopK:              params.TopK,
		TraceId:           traceId,
	})
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func (m *ModelServer) CloseChatCompletionStream() {
	_ = m.streamConn.Close()
}
