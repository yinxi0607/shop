package logic

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shop/order/internal/svc"
	"shop/order/model"
	"shop/order/order"
	"shop/product/product"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type SeckillOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSeckillOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SeckillOrderLogic {
	return &SeckillOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SeckillOrderLogic) SeckillOrder(in *order.SeckillOrderRequest) (*order.SeckillOrderResponse, error) {
	// 获取分布式锁
	lockKey := "seckill:lock:" + in.Pid
	orderKey := "seckill:order:" + in.UserId + ":" + in.Pid

	// 使用SetNX并设置5秒过期时间
	lock, err := l.svcCtx.Redis.SetnxExCtx(l.ctx, lockKey, "1", 5)
	if err != nil {
		return &order.SeckillOrderResponse{Success: false}, err
	}
	if !lock {
		return &order.SeckillOrderResponse{Success: false}, errors.New("高并发抢购中，请稍后重试")
	}
	defer func() {
		_, err = l.svcCtx.Redis.DelCtx(l.ctx, lockKey)
		if err != nil {
			l.Logger.Error("Failed to release lock:", err)
		}
	}()

	// 检查是否重复下单
	isOrderExists, err := l.svcCtx.Redis.ExistsCtx(l.ctx, orderKey)
	if err != nil {
		return &order.SeckillOrderResponse{Success: false}, err
	}
	if isOrderExists {
		return &order.SeckillOrderResponse{Success: false}, errors.New("已抢购过该商品")
	}

	// 查询商品库存
	prodRes, err := l.svcCtx.ProductRpc.GetProductDetail(l.ctx, &product.GetProductDetailRequest{Pid: in.Pid})
	if err != nil || prodRes.Product == nil {
		return &order.SeckillOrderResponse{Success: false}, errors.New("商品不存在")
	}
	if prodRes.Product.Stock < in.Quantity {
		return &order.SeckillOrderResponse{Success: false}, errors.New("库存不足")
	}

	// 扣减库存（使用乐观锁）
	stock := prodRes.Product.Stock - in.Quantity
	updateReq := &product.UpdateProductRequest{
		Pid:   in.Pid,
		Stock: &stock,
	}
	updateRes, err := l.svcCtx.ProductRpc.UpdateProduct(l.ctx, updateReq)
	if err != nil || !updateRes.Success {
		return &order.SeckillOrderResponse{Success: false}, errors.New("库存更新失败")
	}

	// 创建订单（使用事务）
	orderID := uuid.New().String()
	totalPrice := prodRes.Product.Price * float64(in.Quantity)
	err = l.svcCtx.OrderModel.Transact(func(session sqlx.Session) error {
		newOrder := &model.Orders{
			OrderId:    orderID,
			UserId:     in.UserId,
			TotalPrice: totalPrice,
			Status:     "pending",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		if _, err = l.svcCtx.OrderModel.Insert(l.ctx, newOrder); err != nil {
			return err
		}

		orderItem := &model.OrderItems{
			OrderId:   orderID,
			ProductId: in.Pid,
			Quantity:  int64(in.Quantity),
		}
		if _, err = l.svcCtx.OrderItemModel.Insert(l.ctx, orderItem); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return &order.SeckillOrderResponse{Success: false}, err
	}

	// 标记用户已抢购（24小时过期）
	err = l.svcCtx.Redis.SetexCtx(l.ctx, orderKey, "1", 24*60*60) // 改为秒，24小时
	if err != nil {
		return &order.SeckillOrderResponse{Success: false}, err
	}

	return &order.SeckillOrderResponse{
		OrderId: orderID,
		Success: true,
	}, nil
}
