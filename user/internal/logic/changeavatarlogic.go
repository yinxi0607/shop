package logic

import (
	"context"
	"errors"
	"regexp"
	"shop/user/internal/svc"
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

func (l *ChangeAvatarLogic) ChangeAvatar(req *user.ChangeAvatarRequest) (*user.ChangeAvatarResponse, error) {
	// Validate avatar URL (optional)
	if req.NewAvatar != "" {
		if matched, _ := regexp.MatchString(`^https?://[a-zA-Z0-9./_-]+\.(jpg|png|gif)$`, req.NewAvatar); !matched {
			return nil, errors.New("invalid avatar URL")
		}
	}

	// Find user
	u, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, req.UserId)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if u.DeletedAt.Valid {
		return nil, errors.New("user has been deleted")
	}

	// Update avatar
	u.Avatar = req.NewAvatar
	if err = l.svcCtx.UserModel.Update(l.ctx, u); err != nil {
		return nil, err
	}

	return &user.ChangeAvatarResponse{Success: true}, nil
}
