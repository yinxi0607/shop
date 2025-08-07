package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	//Redis redis.RedisConf
	Jwt struct {
		AccessSecret string
		AccessExpire int64
	}
}
