package logic

import (
	"context"
	"errors"
	"shop/user/internal/svc"
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

func (l *ChangeUsernameLogic) ChangeUsername(req *user.ChangeUsernameRequest) (*user.ChangeUsernameResponse, error) {
	// Check if new username exists
	_, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.NewUsername)
	if err == nil {
		return nil, errors.New("new username already exists")
	}

	// Find user
	u, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, req.UserId)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if u.DeletedAt.Valid {
		return nil, errors.New("user has been deleted")
	}

	// Update username
	u.Username = req.NewUsername
	if err = l.svcCtx.UserModel.Update(l.ctx, u); err != nil {
		return nil, err
	}

	return &user.ChangeUsernameResponse{Success: true}, nil
}
