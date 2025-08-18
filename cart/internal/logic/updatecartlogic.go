package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"shop/cart/cart"
	"shop/cart/internal/svc"
	"shop/cart/model"
	"shop/product/product"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCartLogic {
	return &UpdateCartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateCartLogic) UpdateCart(in *cart.UpdateCartRequest) (*cart.UpdateCartResponse, error) {
	// 验证商品
	prodRes, err := l.svcCtx.ProductRpc.GetProductDetail(l.ctx, &product.GetProductDetailRequest{Pid: in.Pid})
	if err != nil || prodRes.Product == nil {
		return &cart.UpdateCartResponse{Success: false}, errors.New("商品不存在")
	}
	if in.Quantity > 0 && prodRes.Product.Stock < in.Quantity {
		return &cart.UpdateCartResponse{Success: false}, errors.New("库存不足")
	}

	// 查询购物车
	cartQuery, err := l.svcCtx.CartModel.FindOneByUserId(l.ctx, in.UserId)
	if err != nil {
		return &cart.UpdateCartResponse{Success: false}, err
	}

	// 更新或软删除购物车项
	now := time.Now()
	if in.Quantity == 0 {
		// 软删除
		if _, err := l.svcCtx.CartItemsModel.SoftDeleteByCartIdAndProductId(l.ctx, cartQuery.CartId, in.Pid); err != nil {
			return &cart.UpdateCartResponse{Success: false}, err
		}
	} else {
		// 更新数量
		item, err := l.svcCtx.CartItemsModel.FindOneByCartIdProductId(l.ctx, cartQuery.CartId, in.Pid)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return &cart.UpdateCartResponse{Success: false}, err
		}
		if item != nil {
			item.Quantity = int64(in.Quantity)
			item.UpdatedAt = now
			if err := l.svcCtx.CartItemsModel.Update(l.ctx, item); err != nil {
				return &cart.UpdateCartResponse{Success: false}, err
			}
		} else {
			newItem := &model.CartItems{
				CartId:    cartQuery.CartId,
				Pid:       in.Pid,
				Quantity:  int64(in.Quantity),
				CreatedAt: now,
				UpdatedAt: now,
			}
			if _, err := l.svcCtx.CartItemsModel.Insert(l.ctx, newItem); err != nil {
				return &cart.UpdateCartResponse{Success: false}, err
			}
		}
	}

	// 更新Redis缓存
	cacheKey := "cart:" + in.UserId
	items, err := l.svcCtx.CartItemsModel.FindByCartId(l.ctx, cartQuery.CartId)
	if err == nil {
		cartItems := make([]*cart.CartItem, len(items))
		for i, item := range items {
			cartItems[i] = &cart.CartItem{Pid: item.Pid, Quantity: int32(item.Quantity)}
		}
		jsonData, err := json.Marshal(&cart.GetCartResponse{Items: cartItems, CartId: cartQuery.CartId})
		if err == nil {
			err = l.svcCtx.Redis.SetexCtx(l.ctx, cacheKey, string(jsonData), 3600)
			if err != nil {
				l.Logger.Error("Failed to Set cart for cache:", err)
			}
		} else {
			l.Logger.Error("Failed to marshal cart for cache:", err)
		}
	}

	return &cart.UpdateCartResponse{Success: true}, nil
}
