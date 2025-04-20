package client

import (
	"context"
	"errors"
	"llm_online_interence/llmgateway/confparser"
	"llm_online_interence/llmgateway/pb"
)

const UserCenterGRPCName = "user_center"

type UserCenterGRPCClient struct {
	conf confparser.GRPCConfigItem
}

func NewUserCenterGRPCClient() (*UserCenterGRPCClient, error) {
	conf, ok := confparser.ProxyConfig.GRPC[UserCenterGRPCName]
	if !ok {
		return nil, errors.New(UserCenterGRPCName + " is not found")
	}
	return &UserCenterGRPCClient{conf: conf}, nil
}

func (u *UserCenterGRPCClient) CheckUserAuth(ctx context.Context, tokenString string) (hasAuth bool, err error) {
	conn, err := newGRPCConn(u.conf)
	if err != nil {
		return false, err
	}
	defer func() { _ = conn.Close() }()

	clt := pb.NewUserCenterServiceClient(conn)
	result, err := clt.CheckAuth(ctx, &pb.UserToken{TokenString: tokenString})
	if err != nil {
		return false, err
	}
	return result.HasAuth, nil
}

func (u *UserCenterGRPCClient) GetUserInfo(ctx context.Context, tokenString string) (userInfo *pb.UserInfo, err error) {
	conn, err := newGRPCConn(u.conf)
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()

	clt := pb.NewUserCenterServiceClient(conn)
	result, err := clt.GetUserInfo(ctx, &pb.UserToken{TokenString: tokenString})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *UserCenterGRPCClient) UpdateUserToken(ctx context.Context, tokenString string) error {
	conn, err := newGRPCConn(u.conf)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	clt := pb.NewUserCenterServiceClient(conn)
	_, err = clt.UpdateUserToken(ctx, &pb.UserToken{TokenString: tokenString})
	return err
}
