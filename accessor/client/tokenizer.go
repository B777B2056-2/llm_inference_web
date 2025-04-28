package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"llm_online_inference/accessor/confparser"
	"llm_online_inference/accessor/pb"
)

type Tokenizer struct {
	addr string
}

func NewTokenizer() *Tokenizer {
	return &Tokenizer{
		addr: fmt.Sprintf("%s:%d", confparser.RPCConfig.Tokenizer.Host, confparser.RPCConfig.Tokenizer.Port),
	}
}

func (t *Tokenizer) Encode(ctx context.Context, prompt string) (inputIds, tokenTypIds []uint64, err error) {
	conn, err := grpc.NewClient(t.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = conn.Close() }()
	clt := pb.NewTokenizerServiceClient(conn)
	traceId, ok := ctx.Value("trace_id").(string)
	if !ok {
		traceId = "unknown"
	}
	resp, err := clt.Tokenizer(ctx, &pb.TokenizerReq{Prompt: prompt, TraceId: traceId})
	if err != nil {
		return nil, nil, err
	}
	return resp.TokenIds, resp.TokenTypeIds, nil
}

func (t *Tokenizer) Decode(ctx context.Context, tokenIds []uint64) (text string, err error) {
	conn, err := grpc.NewClient(t.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer func() { _ = conn.Close() }()
	clt := pb.NewTokenizerServiceClient(conn)
	traceId, ok := ctx.Value("trace_id").(string)
	if !ok {
		traceId = "unknown"
	}
	resp, err := clt.DeTokenizer(ctx, &pb.DeTokenizerReq{TokenIds: tokenIds, TraceId: traceId})
	if err != nil {
		return "", err
	}
	return resp.Text, nil
}
