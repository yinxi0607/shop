package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"net/http"
	"shop/gateway/internal/config"
	"shop/gateway/internal/middleware"
	"shop/product/product"
	"shop/user/user"
)

type ServiceContext struct {
	Config        config.Config
	UserRpc       user.UserRpcClient
	ProductRpc    product.ProductRpcClient
	Redis         *redis.Redis
	JwtMiddleware func(next http.HandlerFunc) http.HandlerFunc
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化用户服务 gRPC 客户端
	userRpc, err := zrpc.NewClient(zrpc.RpcClientConf{
		Target: c.UserRpc.Target,
	})
	if err != nil {
		logx.Errorf("NewServiceContext: failed to initialize UserRpc: %v", err)
		panic(err)
	}

	// 初始化商品服务 gRPC 客户端
	productRpc, err := zrpc.NewClient(zrpc.RpcClientConf{
		Target: c.ProductRpc.Target,
	})
	if err != nil {
		logx.Errorf("NewServiceContext: failed to initialize ProductRpc: %v", err)
		panic(err)
	}

	return &ServiceContext{
		Config:        c,
		UserRpc:       user.NewUserRpcClient(userRpc.Conn()),
		ProductRpc:    product.NewProductRpcClient(productRpc.Conn()),
		Redis:         redis.MustNewRedis(c.Redis),
		JwtMiddleware: middleware.NewJwtMiddleware(c.Jwt.AccessSecret).Handle,
	}
}
