package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"shop/cart/cart"
	"shop/cart/internal/svc"
)

type ClearCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClearCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearCartLogic {
	return &ClearCartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ClearCartLogic) ClearCart(in *cart.ClearCartRequest) (*cart.ClearCartResponse, error) {
	// 查询购物车
	cartQuery, err := l.svcCtx.CartModel.FindOneByUserId(l.ctx, in.UserId)
	if err != nil {
		if err.Error() == "购物车不存在" {
			return &cart.ClearCartResponse{Success: true}, nil
		}
		return &cart.ClearCartResponse{Success: false}, err
	}

	// 软删除购物车项和购物车
	if _, err := l.svcCtx.CartItemsModel.SoftDeleteByCartId(l.ctx, cartQuery.CartId); err != nil {
		return &cart.ClearCartResponse{Success: false}, err
	}
	if _, err := l.svcCtx.CartModel.SoftDeleteByUserId(l.ctx, in.UserId); err != nil {
		return &cart.ClearCartResponse{Success: false}, err
	}

	// 清除Redis缓存
	cacheKey := "cart:" + in.UserId
	if _, err := l.svcCtx.Redis.DelCtx(l.ctx, cacheKey); err != nil {
		l.Logger.Error("Failed to delete cache:", err)
	}

	return &cart.ClearCartResponse{Success: true}, nil
}
