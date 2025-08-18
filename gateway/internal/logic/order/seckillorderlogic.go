package order

import (
	"context"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
	"shop/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type SeckillOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSeckillOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SeckillOrderLogic {
	return &SeckillOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SeckillOrderLogic) SeckillOrder(req *types.SeckillOrderRequest) (resp *types.SeckillOrderResponse, err error) {
	res, err := l.svcCtx.OrderRpc.SeckillOrder(l.ctx, &order.SeckillOrderRequest{
		UserId:   req.UserID,
		Pid:      req.Pid,
		Quantity: req.Quantity,
	})
	if err != nil {
		return nil, err
	}
	return &types.SeckillOrderResponse{
		OrderID: res.OrderId,
		Success: res.Success,
	}, nil
}
