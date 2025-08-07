package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"shop/user/internal/svc"
	"shop/user/user"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginResponse, error) {
	u, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, in.Username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)); err != nil {
		return nil, errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": u.Id,
		"exp":     time.Now().Add(time.Duration(l.svcCtx.Config.Jwt.AccessExpire) * time.Second).Unix(),
	})
	tokenString, err := token.SignedString([]byte(l.svcCtx.Config.Jwt.AccessSecret))
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.Redis.SetexCtx(l.ctx, fmt.Sprintf("session:%d", u.Id), tokenString, int(l.svcCtx.Config.Jwt.AccessExpire))
	if err != nil {
		return nil, err
	}

	return &user.LoginResponse{Token: tokenString}, nil
}
