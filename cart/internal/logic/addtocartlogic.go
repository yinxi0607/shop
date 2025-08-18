package logic

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shop/cart/cart"
	"shop/cart/internal/svc"
	"shop/cart/model"
	"shop/product/product"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddToCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddToCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddToCartLogic {
	return &AddToCartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddToCartLogic) AddToCart(in *cart.AddToCartRequest) (*cart.AddToCartResponse, error) {
	// 验证商品
	prodRes, err := l.svcCtx.ProductRpc.GetProductDetail(l.ctx, &product.GetProductDetailRequest{Pid: in.Pid})
	if err != nil || prodRes.Product == nil {
		return &cart.AddToCartResponse{Success: false}, errors.New("商品不存在")
	}
	if prodRes.Product.Stock < in.Quantity {
		return &cart.AddToCartResponse{Success: false}, errors.New("库存不足")
	}

	// 查询或创建购物车
	cartQuery, err := l.svcCtx.CartModel.FindOneByUserId(l.ctx, in.UserId)
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		return &cart.AddToCartResponse{Success: false}, err
	}
	if cartQuery == nil {
		cartQuery = &model.Carts{
			CartId:    uuid.New().String(),
			UserId:    in.UserId,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if _, err := l.svcCtx.CartModel.Insert(l.ctx, cartQuery); err != nil {
			return &cart.AddToCartResponse{Success: false}, err
		}
		return &cart.AddToCartResponse{Success: true}, nil
	}

	// 添加或更新购物车项
	item, err := l.svcCtx.CartItemsModel.FindOneByCartIdProductId(l.ctx, cartQuery.CartId, in.Pid)
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		return &cart.AddToCartResponse{Success: false}, err
	}
	now := time.Now()
	if item != nil {
		item.Quantity += int64(in.Quantity)
		item.UpdatedAt = now
		err = l.svcCtx.CartItemsModel.Update(l.ctx, item)
		if err != nil {
			return &cart.AddToCartResponse{Success: false}, err
		}
	} else {
		newItem := &model.CartItems{
			CartId:    cartQuery.CartId,
			Pid:       in.Pid,
			Quantity:  int64(in.Quantity),
			CreatedAt: now,
			UpdatedAt: now,
		}
		if _, err = l.svcCtx.CartItemsModel.Insert(l.ctx, newItem); err != nil {
			return &cart.AddToCartResponse{Success: false}, err
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
		jsonData, err := json.Marshal(&cart.GetCartResponse{Items: cartItems})
		if err == nil {
			err = l.svcCtx.Redis.SetexCtx(l.ctx, cacheKey, string(jsonData), 3600)
			if err != nil {
				l.Logger.Error("Failed to Set cart for cache:", err)
			}
		} else {
			l.Logger.Error("Failed to marshal cart for cache:", err)
		}
	}

	return &cart.AddToCartResponse{Success: true}, nil
}
