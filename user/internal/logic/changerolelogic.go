package logic

import (
	"context"
	"errors"

	"shop/user/internal/svc"
	"shop/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangeRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeRoleLogic {
	return &ChangeRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangeRoleLogic) ChangeRole(in *user.ChangeRoleRequest) (*user.ChangeRoleResponse, error) {
	// Find user
	u, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if u.DeletedAt.Valid {
		return nil, errors.New("user has been deleted")
	}
	// Update username
	u.Role = in.Role
	if err = l.svcCtx.UserModel.Update(l.ctx, u); err != nil {
		return nil, err
	}

	return &user.ChangeRoleResponse{Success: true}, nil
}
