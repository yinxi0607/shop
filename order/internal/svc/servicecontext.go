package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"shop/order/internal/config"
	"shop/order/model"
	"shop/product/product"
)

type ServiceContext struct {
	Config         config.Config
	OrderModel     model.OrdersModel
	OrderItemModel model.OrderItemsModel
	Redis          *redis.Redis
	ProductRpc     product.ProductRpcClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化商品服务 gRPC 客户端
	productRpc, err := zrpc.NewClient(c.ProductRpc)
	if err != nil {
		logx.Errorf("NewServiceContext: failed to initialize ProductRpc: %v", err)
		panic(err)
	}
	return &ServiceContext{
		Config:         c,
		OrderModel:     model.NewOrdersModel(sqlx.NewMysql(c.Mysql.DataSource)),
		OrderItemModel: model.NewOrderItemsModel(sqlx.NewMysql(c.Mysql.DataSource)),
		Redis:          redis.MustNewRedis(c.Redis.RedisConf),
		ProductRpc:     product.NewProductRpcClient(productRpc.Conn()),
	}
}
