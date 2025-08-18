package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"shop/cart/cart"
	"shop/cart/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCartLogic {
	return &GetCartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCartLogic) GetCart(in *cart.GetCartRequest) (*cart.GetCartResponse, error) {
	// 尝试从Redis获取缓存
	cacheKey := "cart:" + in.UserId
	cached, err := l.svcCtx.Redis.GetCtx(l.ctx, cacheKey)
	if err == nil && cached != "" {
		var resp cart.GetCartResponse
		if err := json.Unmarshal([]byte(cached), &resp); err == nil {
			return &resp, nil
		}
		l.Logger.Error("Failed to unmarshal cached cart:", err)
	}

	// 查询购物车
	cartQuery, err := l.svcCtx.CartModel.FindOneByUserId(l.ctx, in.UserId)
	if err != nil {
		if err.Error() == "购物车不存在" || errors.Is(err, sql.ErrNoRows) {
			// 购物车不存在，返回空数据
			resp := &cart.GetCartResponse{
				CartId: "",
				Items:  []*cart.CartItem{},
			}
			// 缓存空数据
			jsonData, err := json.Marshal(resp)
			if err == nil {
				err = l.svcCtx.Redis.SetexCtx(l.ctx, cacheKey, string(jsonData), 3600)
				if err != nil {
					l.Logger.Error("Failed to cache empty cart:", err)
				}
			} else {
				l.Logger.Error("Failed to marshal empty cart:", err)
			}
			return resp, nil
		}
		return &cart.GetCartResponse{}, err
	}

	// 查询购物车项
	items, err := l.svcCtx.CartItemsModel.FindByCartId(l.ctx, cartQuery.CartId)
	if err != nil {
		return &cart.GetCartResponse{}, err
	}

	cartItems := make([]*cart.CartItem, len(items))
	for i, item := range items {
		cartItems[i] = &cart.CartItem{Pid: item.Pid, Quantity: int32(item.Quantity)}
	}

	resp := &cart.GetCartResponse{
		CartId: cartQuery.CartId,
		Items:  cartItems,
	}

	// 缓存到Redis
	jsonData, err := json.Marshal(resp)
	if err == nil {
		err = l.svcCtx.Redis.SetexCtx(l.ctx, cacheKey, string(jsonData), 3600)
		if err != nil {
			l.Logger.Error("Failed to cache cart:", err)
		}
	} else {
		l.Logger.Error("Failed to marshal cart:", err)
	}

	return resp, nil
}
