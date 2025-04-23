package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"llm_online_inference/scheduler/confparser"
	"llm_online_inference/scheduler/pb"
)

type Tokenizer struct {
	addr string
}

func NewTokenizer() *Tokenizer {
	return &Tokenizer{
		addr: fmt.Sprintf("%s:%d", confparser.RPCConfig.Tokenizer.Host, confparser.RPCConfig.Tokenizer.Port),
	}
}

func (t *Tokenizer) Do(ctx context.Context, parentMessageID, prompt string) (
	inputIDs, attention_mask []uint32, curTokeCnt uint32, err error) {
	conn, err := grpc.NewClient(t.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, 0, err
	}
	defer func() { _ = conn.Close() }()
	clt := pb.NewTokenizerServiceClient(conn)
	resp, err := clt.Tokenizer(ctx, &pb.TokenizerReq{ParentMessageId: parentMessageID, Prompt: prompt})
	if err != nil {
		return nil, nil, 0, err
	}
	return resp.InputIds, resp.AttentionMask, resp.CurrentTokensCnt, nil
}
