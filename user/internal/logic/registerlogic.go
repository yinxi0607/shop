package logic

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"shop/user/internal/svc"
	"shop/user/model"
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

func (l *RegisterLogic) Register(in *user.RegisterRequest) (*user.RegisterResponse, error) {
	_, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, in.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &model.Users{
		Username:     in.Username,
		PasswordHash: string(passwordHash),
		Email:        in.Email,
	}
	result, err := l.svcCtx.UserModel.Insert(l.ctx, u)
	if err != nil {
		return nil, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &user.RegisterResponse{UserId: userID}, nil
}
