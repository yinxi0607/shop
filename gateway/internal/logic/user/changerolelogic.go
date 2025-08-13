package user

import (
	"context"
	"shop/user/user"

	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeRoleLogic {
	return &ChangeRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeRoleLogic) ChangeRole(req *types.ChangeRoleRequest) (resp *types.ChangeRoleResponse, err error) {
	// 调用用户服务的 gRPC 接口
	_, err = l.svcCtx.UserRpc.ChangeRole(l.ctx, &user.ChangeRoleRequest{
		UserId: req.UserID,
		Role:   req.Role,
	})
	if err != nil {
		logx.Errorf("ChangeRoleLogic: failed to call UserRpc.ChangeRole: %v", err)
		return nil, err
	}

	return &types.ChangeRoleResponse{
		Success: true,
	}, nil
}
