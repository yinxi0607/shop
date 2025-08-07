package logic

import (
	"context"
	"errors"
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

func (l *ChangeUsernameLogic) ChangeUsername(req *types.ChangeUsernameRequest) (*types.ChangeUsernameResponse, error) {
	// Input validation
	if len(req.NewUsername) < 3 || len(req.NewUsername) > 50 {
		return nil, errors.New("new username must be 3-50 characters")
	}

	res, err := l.svcCtx.UserRpc.ChangeUsername(l.ctx, &user.ChangeUsernameRequest{
		UserId:      req.UserID,
		NewUsername: req.NewUsername,
	})
	if err != nil {
		return nil, err
	}
	return &types.ChangeUsernameResponse{Success: res.Success}, nil
}
