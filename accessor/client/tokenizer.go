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

func (t *Tokenizer) Encode(ctx context.Context, prompt string) (input_ids, token_typ_ids []uint64, err error) {
	conn, err := grpc.NewClient(t.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = conn.Close() }()
	clt := pb.NewTokenizerServiceClient(conn)
	resp, err := clt.Tokenizer(ctx, &pb.TokenizerReq{Prompt: prompt})
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
	resp, err := clt.DeTokenizer(ctx, &pb.DeTokenizerReq{TokenIds: tokenIds})
	if err != nil {
		return "", err
	}
	return resp.Text, nil
}
