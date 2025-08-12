package logic

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"shop/user/internal/svc"
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

func (l *ChangePasswordLogic) ChangePassword(req *user.ChangePasswordRequest) (*user.ChangePasswordResponse, error) {
	u, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, req.UserId)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if u.DeletedAt.Valid {
		return nil, errors.New("user has been deleted")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.OldPassword)); err != nil {
		return nil, errors.New("invalid old password")
	}

	// Generate new password hash
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Update password
	u.PasswordHash = string(newPasswordHash)
	if err = l.svcCtx.UserModel.Update(l.ctx, u); err != nil {
		return nil, err
	}

	// Invalidate old session
	_, err = l.svcCtx.Redis.DelCtx(l.ctx, fmt.Sprintf("session:%d", u.Id))
	if err != nil {
		return nil, err
	}

	return &user.ChangePasswordResponse{Success: true}, nil
}
