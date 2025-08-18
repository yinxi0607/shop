package logic

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shop/cart/cart"
	"shop/order/internal/svc"
	"shop/order/model"
	"shop/order/order"
	"shop/product/product"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
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
	var items []*order.OrderItem

	// 从购物车获取items
	if in.UseCart {
		cartRes, err := l.svcCtx.CartRpc.GetCart(l.ctx, &cart.GetCartRequest{UserId: in.UserId})
		if err != nil {
			return &order.CreateOrderResponse{}, err
		}
		if len(cartRes.Items) == 0 {
			return &order.CreateOrderResponse{}, errors.New("购物车为空")
		}
		//items = cartRes.Items
		for _, item := range cartRes.Items {
			items = append(items, &order.OrderItem{
				Pid:      item.Pid,
				Quantity: item.Quantity,
			})
		}
	} else {
		return &order.CreateOrderResponse{}, errors.New("非购物车下单暂不支持")
	}

	// 计算总价并验证库存
	var totalPrice float64
	for _, item := range items {
		prodRes, err := l.svcCtx.ProductRpc.GetProductDetail(l.ctx, &product.GetProductDetailRequest{Pid: item.Pid})
		if err != nil || prodRes.Product == nil {
			return &order.CreateOrderResponse{}, errors.New("商品 " + item.Pid + " 不存在")
		}
		if prodRes.Product.Stock < item.Quantity {
			return &order.CreateOrderResponse{}, errors.New("商品 " + item.Pid + " 库存不足")
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

		for _, item := range items {
			orderItem := &model.OrderItems{
				OrderId:   orderID,
				ProductId: item.Pid,
				Quantity:  int64(item.Quantity),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if _, err := l.svcCtx.OrderItemModel.Insert(l.ctx, orderItem); err != nil {
				return err
			}
		}

		// 扣减库存
		for _, item := range items {
			prodRes, err := l.svcCtx.ProductRpc.GetProductDetail(l.ctx, &product.GetProductDetailRequest{Pid: item.Pid})
			if err != nil {
				return err
			}
			stock := prodRes.Product.Stock - item.Quantity
			updateReq := &product.UpdateProductRequest{
				Pid:   item.Pid,
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

	// 清空购物车（软删除）
	if in.UseCart {
		if _, err := l.svcCtx.CartRpc.ClearCart(l.ctx, &cart.ClearCartRequest{UserId: in.UserId}); err != nil {
			l.Logger.Error("Failed to clear cart:", err)
		}
	}

	// 清除订单详情缓存
	cacheKey := "order:detail:" + orderID
	if _, err := l.svcCtx.Redis.DelCtx(l.ctx, cacheKey); err != nil {
		l.Logger.Error("Failed to delete order cache:", err)
	}

	return &order.CreateOrderResponse{OrderId: orderID}, nil
}
