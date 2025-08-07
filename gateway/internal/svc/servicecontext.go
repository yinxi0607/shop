package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"shop/gateway/internal/config"
	"shop/user/user"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc user.UserRpcClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := zrpc.MustNewClient(c.UserRpc).Conn() // 获取 *grpc.ClientConn
	return &ServiceContext{
		Config:  c,
		UserRpc: user.NewUserRpcClient(conn), // 传递 *grpc.ClientConn
	}
}
