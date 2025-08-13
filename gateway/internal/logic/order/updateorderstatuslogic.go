package order

import (
	"context"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
	"shop/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrderStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateOrderStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrderStatusLogic {
	return &UpdateOrderStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateOrderStatusLogic) UpdateOrderStatus(req *types.UpdateOrderStatusRequest) (resp *types.UpdateOrderStatusResponse, err error) {
	res, err := l.svcCtx.OrderRpc.UpdateOrderStatus(l.ctx, &order.UpdateOrderStatusRequest{
		OrderId: req.OrderID,
		Status:  req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &types.UpdateOrderStatusResponse{Success: res.Success}, nil
}
