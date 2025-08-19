package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	OrderRpc zrpc.RpcClientConf
	Kafka    struct {
		Brokers []string
		Topic   string
	}
	WechatPay struct {
		AppId      string
		MchId      string
		SerialNo   string
		ApiV3Key   string
		PrivateKey string
	}
	Alipay struct {
		AppId      string
		PrivateKey string
		PublicKey  string
		NotifyUrl  string
		IsProd     bool
	}
}
