package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"shop/order/internal/svc"
	"shop/order/order"
)

type GetOrderDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderDetailLogic {
	return &GetOrderDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOrderDetailLogic) GetOrderDetail(in *order.GetOrderDetailRequest) (*order.GetOrderDetailResponse, error) {
	// 定义缓存键
	cacheKey := fmt.Sprintf("order:detail:%s:%s", in.UserId, in.OrderId)

	// 尝试从Redis获取缓存
	cached, err := l.svcCtx.Redis.GetCtx(l.ctx, cacheKey)
	if err == nil && cached != "" {
		var resp order.GetOrderDetailResponse
		if err := json.Unmarshal([]byte(cached), &resp); err == nil {
			return &resp, nil
		}
		l.Logger.Error("Failed to unmarshal cached order:", err)
	}

	// 查询订单
	ord, err := l.svcCtx.OrderModel.FindOneByOrderId(l.ctx, in.OrderId, in.UserId)
	if err != nil {
		return &order.GetOrderDetailResponse{}, errors.New("订单不存在")
	}

	// 查询订单项
	items, err := l.svcCtx.OrderItemModel.FindByOrderId(l.ctx, in.OrderId)
	if err != nil {
		return &order.GetOrderDetailResponse{}, err
	}

	orderItems := make([]*order.OrderItem, len(items))
	for i, item := range items {
		orderItems[i] = &order.OrderItem{
			ProductId: item.ProductId,
			Quantity:  int32(item.Quantity),
		}
	}

	resp := &order.GetOrderDetailResponse{
		Order: &order.Order{
			OrderId:    ord.OrderId,
			UserId:     ord.UserId,
			Items:      orderItems,
			TotalPrice: ord.TotalPrice,
			Status:     ord.Status,
			CreatedAt:  ord.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  ord.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}

	// 缓存到Redis（1小时过期）
	jsonData, err := json.Marshal(resp)
	if err == nil {
		err = l.svcCtx.Redis.SetexCtx(l.ctx, cacheKey, string(jsonData), 3600)
		if err != nil {
			l.Logger.Error("Failed to cache order detail:", err)
		}
	} else {
		l.Logger.Error("Failed to marshal order detail:", err)
	}

	return resp, nil
}
