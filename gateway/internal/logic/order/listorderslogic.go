package order

import (
	"context"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
	"shop/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOrdersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOrdersLogic {
	return &ListOrdersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListOrdersLogic) ListOrders(req *types.ListOrdersRequest) (resp *types.ListOrdersResponse, err error) {
	res, err := l.svcCtx.OrderRpc.ListOrders(l.ctx, &order.ListOrdersRequest{
		UserId:   req.UserID,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	orders := make([]types.Order, len(res.Orders))
	for i, ord := range res.Orders {
		items := make([]types.OrderItem, len(ord.Items))
		for j, item := range ord.Items {
			items[j] = types.OrderItem{ProductID: item.ProductId, Quantity: item.Quantity}
		}
		orders[i] = types.Order{
			OrderID:    ord.OrderId,
			UserID:     ord.UserId,
			Items:      items,
			TotalPrice: ord.TotalPrice,
			Status:     ord.Status,
			CreatedAt:  ord.CreatedAt,
			UpdatedAt:  ord.UpdatedAt,
		}
	}
	return &types.ListOrdersResponse{
		Orders: orders,
		Total:  res.Total,
	}, nil
}
