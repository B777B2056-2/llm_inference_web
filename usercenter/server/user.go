package server

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"llm_inference_web/usercenter/pb"
	"llm_inference_web/usercenter/services"
	"llm_inference_web/usercenter/utils"
)

type UserCenterServer struct {
	pb.UnimplementedUserCenterServiceServer
}

func (u *UserCenterServer) CheckAuth(ctx context.Context, param *pb.UserToken) (*pb.AuthCheckResult, error) {
	// 校验token是否有效
	_, err := utils.ValidateToken(param.TokenString)
	if err != nil {
		return nil, err
	}

	// 检查用户是否已登录
	hasAuth, err := services.CheckUserAuth(ctx, param.TokenString)
	if err != nil {
		return nil, err
	}
	return &pb.AuthCheckResult{
		HasAuth: hasAuth,
	}, nil
}

func (u *UserCenterServer) GetUserInfo(ctx context.Context, param *pb.UserToken) (*pb.UserInfo, error) {
	// 校验token是否有效
	_, err := utils.ValidateToken(param.TokenString)
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	userInfo, err := services.GetUserInfo(ctx, param.TokenString)
	if err != nil {
		return nil, err
	}
	return &pb.UserInfo{
		Id:   uint32(userInfo.UserID),
		Name: userInfo.Username,
	}, nil
}

func (u *UserCenterServer) UpdateUserToken(ctx context.Context, param *pb.UserToken) (*emptypb.Empty, error) {
	// 校验token是否有效
	_, err := utils.ValidateToken(param.TokenString)
	if err != nil {
		return nil, err
	}

	return nil, services.UpdateUserToken(ctx, param.TokenString)
}
