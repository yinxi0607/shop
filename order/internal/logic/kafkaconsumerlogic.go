package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"shop/order/internal/svc"

	"github.com/IBM/sarama"
	"github.com/zeromicro/go-zero/core/logx"
)

type Notification struct {
	Type          string  `json:"type"` // payment, refund
	PaymentId     string  `json:"payment_id"`
	RefundId      string  `json:"refund_id,omitempty"`
	OrderId       string  `json:"order_id"`
	Status        string  `json:"status"`
	Amount        float64 `json:"amount"`
	TransactionId string  `json:"transaction_id"`
}

type KafkaConsumerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewKafkaConsumerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KafkaConsumerLogic {
	return &KafkaConsumerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *KafkaConsumerLogic) Consume() error {
	consumer, err := sarama.NewConsumer(l.svcCtx.Config.Kafka.Brokers, nil)
	if err != nil {
		return fmt.Errorf("创建 Kafka 消费者失败: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(l.svcCtx.Config.Kafka.Topic, 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("消费分区失败: %v", err)
	}
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
		var notification Notification
		if err := json.Unmarshal(msg.Value, &notification); err != nil {
			l.Logger.Errorf("解析通知失败: %v", err)
			continue
		}
		userId := ""
		switch notification.Type {
		case "payment":
			if notification.Status == "completed" {
				err = l.updateOrderStatus(notification.OrderId, userId, "paid")
				if err != nil {
					l.Logger.Errorf("更新支付订单状态失败: %v", err)
				}
			}
		case "refund":
			if notification.Status == "completed" {
				err = l.updateOrderStatus(notification.OrderId, userId, "refunded")
				if err != nil {
					l.Logger.Errorf("更新退款订单状态失败: %v", err)
				}
			}
		default:
			l.Logger.Errorf("未知通知类型: %s", notification.Type)
		}
	}
	return nil
}

func (l *KafkaConsumerLogic) updateOrderStatus(orderId, userId, status string) error {
	// 查询订单
	order, err := l.svcCtx.OrderModel.FindOneByOrderId(l.ctx, orderId, userId)
	if err != nil || order == nil {
		return fmt.Errorf("订单不存在: %s", orderId)
	}

	// 更新订单状态
	err = l.svcCtx.OrderModel.UpdateStates(l.ctx, orderId, userId, status)
	if err != nil {
		return fmt.Errorf("更新订单状态失败: %v", err)
	}

	// 更新缓存
	cacheKey := "order:detail:" + orderId
	err = l.svcCtx.Redis.SetexCtx(l.ctx, cacheKey, status, 3600)
	if err != nil {
		l.Logger.Errorf("更新订单缓存失败: %v", err)
	}

	return nil
}
