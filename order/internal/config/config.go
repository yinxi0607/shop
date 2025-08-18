package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	ProductRpc zrpc.RpcClientConf
	CartRpc    zrpc.RpcClientConf
	Mysql      struct {
		DataSource string
	}
}
