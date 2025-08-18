package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"shop/cart/internal/config"
	"shop/cart/model"
	"shop/order/order"
	"shop/product/product"
)

type ServiceContext struct {
	Config         config.Config
	CartModel      model.CartsModel
	CartItemsModel model.CartItemsModel
	Redis          *redis.Redis
	ProductRpc     product.ProductRpcClient
	OrderRpc       order.OrderRpcClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化商品服务 gRPC 客户端
	productRpc, err := zrpc.NewClient(c.ProductRpc)
	if err != nil {
		logx.Errorf("NewServiceContext: failed to initialize ProductRpc: %v", err)
		panic(err)
	}

	// 初始化订单服务 gRPC 客户端
	orderRpc, err := zrpc.NewClient(c.OrderRpc)
	if err != nil {
		logx.Errorf("NewServiceContext: failed to initialize OrderRpc: %v", err)
		panic(err)
	}
	return &ServiceContext{
		Config:         c,
		CartModel:      model.NewCartsModel(sqlx.NewMysql(c.Mysql.DataSource)),
		CartItemsModel: model.NewCartItemsModel(sqlx.NewMysql(c.Mysql.DataSource)),
		Redis:          redis.MustNewRedis(c.Redis.RedisConf),
		ProductRpc:     product.NewProductRpcClient(productRpc.Conn()),
		OrderRpc:       order.NewOrderRpcClient(orderRpc.Conn()),
	}
}
