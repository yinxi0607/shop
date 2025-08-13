package order

import (
	"context"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
	"shop/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderRequest) (resp *types.CreateOrderResponse, err error) {
	items := make([]*order.OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = &order.OrderItem{ProductId: item.ProductID, Quantity: item.Quantity}
	}
	res, err := l.svcCtx.OrderRpc.CreateOrder(l.ctx, &order.CreateOrderRequest{
		UserId: req.UserID,
		Items:  items,
	})
	if err != nil {
		return nil, err
	}
	return &types.CreateOrderResponse{OrderID: res.OrderId}, nil
}
