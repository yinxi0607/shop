package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	ProductRpc zrpc.RpcClientConf
	OrderRpc   zrpc.RpcClientConf
	Mysql      struct {
		DataSource string
	}
}
