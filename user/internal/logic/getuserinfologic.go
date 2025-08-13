package logic

import (
	"context"
	"errors"
	"shop/user/internal/svc"
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

func (l *GetUserInfoLogic) GetUserInfo(req *user.GetUserInfoRequest) (*user.GetUserInfoResponse, error) {
	u, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, req.UserId)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if u.DeletedAt.Valid {
		return nil, errors.New("user has been deleted")
	}

	address := ""
	if u.Address.Valid {
		address = u.Address.String
	}

	return &user.GetUserInfoResponse{
		User: &user.UserInfo{
			Id:       u.Id,
			Username: u.Username,
			Email:    u.Email,
			Avatar:   u.Avatar,
			Bio:      u.Bio,
			Address:  address,
			UserId:   u.UserId,
			CreateAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
