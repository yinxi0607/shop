package user

import (
	"context"
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

func (l *ChangeAvatarLogic) ChangeAvatar(req *types.ChangeAvatarRequest) (resp *types.ChangeAvatarResponse, err error) {
	// 调用用户服务的 gRPC 接口
	_, err = l.svcCtx.UserRpc.ChangeAvatar(l.ctx, &user.ChangeAvatarRequest{
		UserId:    req.UserID,
		NewAvatar: req.NewAvatar,
	})
	if err != nil {
		logx.Errorf("ChangeAvatarLogic: failed to call UserRpc.ChangeAvatar: %v", err)
		return nil, err
	}

	return &types.ChangeAvatarResponse{
		Success: true,
	}, nil
}
