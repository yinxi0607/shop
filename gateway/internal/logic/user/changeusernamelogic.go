package user

import (
	"context"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
	"shop/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeUsernameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeUsernameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeUsernameLogic {
	return &ChangeUsernameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeUsernameLogic) ChangeUsername(req *types.ChangeUsernameRequest) (resp *types.ChangeUsernameResponse, err error) {
	// 调用用户服务的 gRPC 接口
	_, err = l.svcCtx.UserRpc.ChangeUsername(l.ctx, &user.ChangeUsernameRequest{
		UserId:      req.UserID,
		NewUsername: req.NewUsername,
	})
	if err != nil {
		logx.Errorf("ChangeUsernameLogic: failed to call UserRpc.ChangeUsername: %v", err)
		return nil, err
	}

	return &types.ChangeUsernameResponse{
		Success: true,
	}, nil
}
