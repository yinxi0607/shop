package user

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
	"shop/user/user"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// 调用用户服务的 gRPC 接口
	loginResp, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		logx.Errorf("LoginLogic: failed to call UserRpc.Login: %v", err)
		return nil, err
	}

	return &types.LoginResponse{
		Token: loginResp.Token,
	}, nil
}
