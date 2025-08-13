package logic

import (
	"context"
	"errors"
	"shop/order/internal/svc"
	"shop/order/order"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrderStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOrderStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrderStatusLogic {
	return &UpdateOrderStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateOrderStatusLogic) UpdateOrderStatus(in *order.UpdateOrderStatusRequest) (*order.UpdateOrderStatusResponse, error) {
	// 验证订单是否存在
	ord, err := l.svcCtx.OrderModel.FindOneByOrderId(l.ctx, in.OrderId)
	if err != nil {
		return &order.UpdateOrderStatusResponse{Success: false}, errors.New("订单不存在")
	}

	// 验证状态合法性
	validStatuses := map[string]bool{"pending": true, "paid": true, "shipped": true, "cancelled": true}
	if !validStatuses[in.Status] {
		return &order.UpdateOrderStatusResponse{Success: false}, errors.New("无效的订单状态")
	}

	// 更新订单状态
	ord.Status = in.Status
	ord.UpdatedAt = time.Now()
	if err = l.svcCtx.OrderModel.Update(l.ctx, ord); err != nil {
		return &order.UpdateOrderStatusResponse{Success: false}, err
	}

	// 清除Redis缓存
	cacheKey := "order:detail:" + in.OrderId
	if _, err := l.svcCtx.Redis.DelCtx(l.ctx, cacheKey); err != nil {
		l.Logger.Error("Failed to delete cache:", err)
	}

	return &order.UpdateOrderStatusResponse{Success: true}, nil
}
