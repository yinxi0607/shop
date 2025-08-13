package user

import (
	"context"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
	"shop/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo(req *types.GetUserInfoRequest) (resp *types.GetUserInfoResponse, err error) {
	// 调用用户服务的 gRPC 接口
	userResp, err := l.svcCtx.UserRpc.GetUserInfo(l.ctx, &user.GetUserInfoRequest{
		UserId: req.UserID,
	})
	if err != nil {
		logx.Errorf("GetUserInfoLogic: failed to call UserRpc.GetUserInfo: %v", err)
		return nil, err
	}

	return &types.GetUserInfoResponse{
		UserID:    userResp.User.UserId,
		Username:  userResp.User.Username,
		Email:     userResp.User.Email,
		Avatar:    userResp.User.Avatar,
		Bio:       userResp.User.Bio,
		Address:   userResp.User.Address,
		CreatedAt: userResp.User.CreateAt,
	}, nil
}
