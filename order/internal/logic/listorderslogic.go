package logic

import (
	"context"
	"shop/order/internal/svc"
	"shop/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOrdersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOrdersLogic {
	return &ListOrdersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListOrdersLogic) ListOrders(in *order.ListOrdersRequest) (*order.ListOrdersResponse, error) {
	// 查询订单总数
	total, err := l.svcCtx.OrderModel.CountByUserId(l.ctx, in.UserId)
	if err != nil {
		return &order.ListOrdersResponse{}, err
	}

	// 分页查询订单
	orders, err := l.svcCtx.OrderModel.FindByUserId(l.ctx, in.UserId, in.Page, in.PageSize)
	if err != nil {
		return &order.ListOrdersResponse{}, err
	}

	resOrders := make([]*order.Order, len(orders))
	for i, ord := range orders {
		// 查询订单项
		items, err := l.svcCtx.OrderItemModel.FindByOrderId(l.ctx, ord.OrderId)
		if err != nil {
			return &order.ListOrdersResponse{}, err
		}
		orderItems := make([]*order.OrderItem, len(items))
		for j, item := range items {
			orderItems[j] = &order.OrderItem{
				Pid:      item.ProductId,
				Quantity: int32(item.Quantity),
			}
		}
		resOrders[i] = &order.Order{
			OrderId:    ord.OrderId,
			UserId:     ord.UserId,
			Items:      orderItems,
			TotalPrice: ord.TotalPrice,
			Status:     ord.Status,
			CreatedAt:  ord.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  ord.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &order.ListOrdersResponse{
		Orders: resOrders,
		Total:  int32(total),
	}, nil
}
