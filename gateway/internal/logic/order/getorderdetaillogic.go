package order

import (
	"context"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
	"shop/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderDetailLogic {
	return &GetOrderDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderDetailLogic) GetOrderDetail(req *types.GetOrderDetailRequest) (resp *types.GetOrderDetailResponse, err error) {
	res, err := l.svcCtx.OrderRpc.GetOrderDetail(l.ctx, &order.GetOrderDetailRequest{OrderId: req.OrderID, UserId: req.UserID})
	if err != nil {
		return nil, err
	}
	items := make([]types.OrderItem, len(res.Order.Items))
	for i, item := range res.Order.Items {
		items[i] = types.OrderItem{Pid: item.Pid, Quantity: item.Quantity}
	}
	return &types.GetOrderDetailResponse{
		Order: types.Order{
			OrderID:    res.Order.OrderId,
			UserID:     res.Order.UserId,
			Items:      items,
			TotalPrice: res.Order.TotalPrice,
			Status:     res.Order.Status,
			CreatedAt:  res.Order.CreatedAt,
			UpdatedAt:  res.Order.UpdatedAt,
		},
	}, nil
}
