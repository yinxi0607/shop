package user

import (
	"context"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
	"shop/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	// 调用用户服务的 gRPC 接口
	userResp, err := l.svcCtx.UserRpc.Register(l.ctx, &user.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Avatar:   req.Avatar,
		Bio:      req.Bio,
		Address:  req.Address,
	})
	if err != nil {
		logx.Errorf("RegisterLogic: failed to call UserRpc.Register: %v", err)
		return nil, err
	}

	return &types.RegisterResponse{
		UserID: userResp.UserId,
	}, nil
}
