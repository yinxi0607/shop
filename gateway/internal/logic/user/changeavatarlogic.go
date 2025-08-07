package logic

import (
	"context"
	"errors"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
	"shop/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeAvatarLogic {
	return &ChangeAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeAvatarLogic) ChangeAvatar(req *types.ChangeAvatarRequest) (*types.ChangeAvatarResponse, error) {
	// Input validation
	if len(req.NewAvatar) > 255 {
		return nil, errors.New("avatar URL too long")
	}

	res, err := l.svcCtx.UserRpc.ChangeAvatar(l.ctx, &user.ChangeAvatarRequest{
		UserId:    req.UserID,
		NewAvatar: req.NewAvatar,
	})
	if err != nil {
		return nil, err
	}
	return &types.ChangeAvatarResponse{Success: res.Success}, nil
}
