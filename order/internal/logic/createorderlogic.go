package logic

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"shop/order/internal/svc"
	"shop/order/model"
	"shop/order/order"
	"shop/product/product"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOrderLogic) CreateOrder(in *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	// 计算总价并验证库存
	var totalPrice float64
	for _, item := range in.Items {
		prodRes, err := l.svcCtx.ProductRpc.GetProductDetail(l.ctx, &product.GetProductDetailRequest{Pid: item.ProductId})
		if err != nil || prodRes.Product == nil {
			return &order.CreateOrderResponse{}, errors.New("商品 " + item.ProductId + " 不存在")
		}
		if prodRes.Product.Stock < item.Quantity {
			return &order.CreateOrderResponse{}, errors.New("商品 " + item.ProductId + " 库存不足")
		}
		totalPrice += prodRes.Product.Price * float64(item.Quantity)
	}

	// 创建订单（使用事务）
	orderID := uuid.New().String()
	err := l.svcCtx.OrderModel.Transact(func(session sqlx.Session) error {
		newOrder := &model.Orders{
			OrderId:    orderID,
			UserId:     in.UserId,
			TotalPrice: totalPrice,
			Status:     "pending",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		if _, err := l.svcCtx.OrderModel.Insert(l.ctx, newOrder); err != nil {
			return err
		}

		for _, item := range in.Items {
			orderItem := &model.OrderItems{
				OrderId:   orderID,
				ProductId: item.ProductId,
				Quantity:  int64(item.Quantity),
			}
			if _, err := l.svcCtx.OrderItemModel.Insert(l.ctx, orderItem); err != nil {
				return err
			}
		}

		// 扣减库存
		for _, item := range in.Items {
			prodRes, err := l.svcCtx.ProductRpc.GetProductDetail(l.ctx, &product.GetProductDetailRequest{Pid: item.ProductId})
			if err != nil {
				return err
			}
			stock := prodRes.Product.Stock - item.Quantity
			updateReq := &product.UpdateProductRequest{
				Pid:   item.ProductId,
				Stock: &stock,
			}
			if updateRes, err := l.svcCtx.ProductRpc.UpdateProduct(l.ctx, updateReq); err != nil || !updateRes.Success {
				return errors.New("库存更新失败")
			}
		}

		return nil
	})
	if err != nil {
		return &order.CreateOrderResponse{}, err
	}

	return &order.CreateOrderResponse{OrderId: orderID}, nil
}
