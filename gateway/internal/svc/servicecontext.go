package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"net/http"
	"shop/gateway/internal/config"
	"shop/gateway/internal/middleware"
	"shop/order/order"
	"shop/product/product"
	"shop/user/user"
)

type ServiceContext struct {
	Config        config.Config
	UserRpc       user.UserRpcClient
	ProductRpc    product.ProductRpcClient
	OrderRpc      order.OrderRpcClient
	Redis         *redis.Redis
	JwtMiddleware func(next http.HandlerFunc) http.HandlerFunc
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化用户服务 gRPC 客户端
	userRpc, err := zrpc.NewClient(c.UserRpc)
	if err != nil {
		logx.Errorf("NewServiceContext: failed to initialize UserRpc: %v", err)
		panic(err)
	}

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
		Config:        c,
		UserRpc:       user.NewUserRpcClient(userRpc.Conn()),
		ProductRpc:    product.NewProductRpcClient(productRpc.Conn()),
		OrderRpc:      order.NewOrderRpcClient(orderRpc.Conn()),
		Redis:         redis.MustNewRedis(c.Redis),
		JwtMiddleware: middleware.NewJwtMiddleware(c.Jwt.AccessSecret).Handle,
	}
}
