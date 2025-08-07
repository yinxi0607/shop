package logic

import (
	"context"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
	"shop/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordRequest) (*types.ChangePasswordResponse, error) {
	res, err := l.svcCtx.UserRpc.ChangePassword(l.ctx, &user.ChangePasswordRequest{
		UserId:      req.UserID,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		return nil, err
	}
	return &types.ChangePasswordResponse{Success: res.Success}, nil
}
