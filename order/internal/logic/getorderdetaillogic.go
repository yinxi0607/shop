package logic

import (
	"context"
	"errors"
	"shop/order/internal/svc"
	"shop/order/order"

	"github.com/zeromicro/go-zero/core/logx"
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
	// 查询订单
	ord, err := l.svcCtx.OrderModel.FindOneByOrderId(l.ctx, in.OrderId)
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

	return &order.GetOrderDetailResponse{
		Order: &order.Order{
			OrderId:    ord.OrderId,
			UserId:     ord.UserId,
			Items:      orderItems,
			TotalPrice: ord.TotalPrice,
			Status:     ord.Status,
			CreatedAt:  ord.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  ord.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
