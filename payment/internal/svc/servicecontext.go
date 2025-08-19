package svc

import (
	"github.com/IBM/sarama"
	"github.com/go-pay/gopay/alipay/v3"
	"github.com/go-pay/gopay/wechat/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"shop/order/order"
	"shop/payment/internal/config"
	"shop/payment/model"
)

type ServiceContext struct {
	Config        config.Config
	PaymentModel  model.PaymentsModel
	RefundModel   model.RefundsModel
	Redis         *redis.Redis
	OrderRpc      order.OrderRpcClient
	WechatClient  *wechat.ClientV3
	AlipayClient  *alipay.ClientV3
	KafkaProducer sarama.SyncProducer
}

func NewServiceContext(c config.Config) *ServiceContext {
	producer, err := sarama.NewSyncProducer(c.Kafka.Brokers, nil)
	if err != nil {
		panic(err)
	}
	aliPayV3, err := alipay.NewClientV3(c.Alipay.AppId, c.Alipay.PrivateKey, c.Alipay.IsProd)
	if err != nil {
		return nil
	}
	wechatV3, err := wechat.NewClientV3(c.WechatPay.MchId, c.WechatPay.SerialNo, c.WechatPay.ApiV3Key, c.WechatPay.PrivateKey)
	if err != nil {
		return nil
	}
	// 初始化订单服务 gRPC 客户端
	orderRpc, err := zrpc.NewClient(c.OrderRpc)
	if err != nil {
		logx.Errorf("NewServiceContext: failed to initialize OrderRpc: %v", err)
		panic(err)
	}
	return &ServiceContext{
		Config:        c,
		PaymentModel:  model.NewPaymentsModel(sqlx.NewMysql(c.Mysql.DataSource)),
		RefundModel:   model.NewRefundsModel(sqlx.NewMysql(c.Mysql.DataSource)),
		Redis:         redis.New(c.Redis.Host, redis.WithPass(c.Redis.Pass)),
		OrderRpc:      order.NewOrderRpcClient(orderRpc.Conn()),
		WechatClient:  wechatV3,
		AlipayClient:  aliPayV3,
		KafkaProducer: producer,
	}
}
