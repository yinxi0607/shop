package logic

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

func (l *GetUserInfoLogic) GetUserInfo(req *types.GetUserInfoRequest) (*types.GetUserInfoResponse, error) {
	res, err := l.svcCtx.UserRpc.GetUserInfo(l.ctx, &user.GetUserInfoRequest{
		UserId: req.UserID,
	})
	if err != nil {
		return nil, err
	}
	return &types.GetUserInfoResponse{
		User: types.UserInfo{
			ID:       res.User.Id,
			Username: res.User.Username,
			Email:    res.User.Email,
			Avatar:   res.User.Avatar,
			Bio:      res.User.Bio,
			Address:  res.User.Address,
		},
	}, nil
}
