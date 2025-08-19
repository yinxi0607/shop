package logic

import (
	"context"
	"fmt"
	"shop/order/order"
	"shop/payment/internal/svc"
	"shop/payment/payment"

	"github.com/IBM/sarama"
	"github.com/zeromicro/go-zero/core/logx"
)

type PaymentCallbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPaymentCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentCallbackLogic {
	return &PaymentCallbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PaymentCallbackLogic) PaymentCallback(req *payment.PaymentCallbackRequest) (resp *payment.PaymentCallbackResponse, err error) {

	// 防止重放攻击
	replayKey := "callback:replay:" + req.TransactionId
	exists, err := l.svcCtx.Redis.ExistsCtx(l.ctx, replayKey)
	if err != nil || exists {
		return nil, fmt.Errorf("重复的回调请求")
	}
	err = l.svcCtx.Redis.SetexCtx(l.ctx, replayKey, "1", 3600*24)
	if err != nil {
		return nil, fmt.Errorf("记录回调失败: %v", err)
	}

	// 查询支付记录
	paymentQuery, err := l.svcCtx.PaymentModel.FindOneByPaymentId(l.ctx, req.PaymentId)
	if err != nil || paymentQuery == nil {
		return nil, fmt.Errorf("支付记录不存在")
	}
	if paymentQuery.Status != "pending" {
		return nil, fmt.Errorf("支付状态已更新: %s", paymentQuery.Status)
	}

	// 更新支付状态
	err = l.svcCtx.PaymentModel.UpdateStatus(l.ctx, req.PaymentId, req.Status, req.TransactionId)
	if err != nil {
		return nil, fmt.Errorf("更新支付状态失败: %v", err)
	}

	// 更新订单状态
	if req.Status == "completed" {
		_, err = l.svcCtx.OrderRpc.UpdateOrderStatus(l.ctx, &order.UpdateOrderStatusRequest{
			OrderId: paymentQuery.OrderId,
			Status:  "paid",
		})
		if err != nil {
			return nil, fmt.Errorf("更新订单状态失败: %v", err)
		}
	}

	// 更新缓存
	err = l.svcCtx.Redis.SetexCtx(l.ctx, "payment:"+paymentQuery.OrderId, req.Status, 3600)
	if err != nil {
		return nil, fmt.Errorf("更新缓存失败: %v", err)
	}

	// 发送 Kafka 通知
	msg := &struct {
		Type          string  `json:"type"`
		PaymentId     string  `json:"payment_id"`
		OrderId       string  `json:"order_id"`
		Status        string  `json:"status"`
		Amount        float64 `json:"amount"`
		TransactionId string  `json:"transaction_id"`
	}{
		Type:          "payment",
		PaymentId:     req.PaymentId,
		OrderId:       paymentQuery.OrderId,
		Status:        req.Status,
		Amount:        paymentQuery.Amount,
		TransactionId: req.TransactionId,
	}
	_, _, err = l.svcCtx.KafkaProducer.SendMessage(&sarama.ProducerMessage{
		Topic: l.svcCtx.Config.Kafka.Topic,
		Value: sarama.StringEncoder(fmt.Sprintf(`{"type":"%s","payment_id":"%s","order_id":"%s","status":"%s","amount":%.2f,"transaction_id":"%s"}`,
			msg.Type, msg.PaymentId, msg.OrderId, msg.Status, msg.Amount, msg.TransactionId)),
	})
	if err != nil {
		return nil, fmt.Errorf("发送支付通知失败: %v", err)
	}

	return &payment.PaymentCallbackResponse{Success: true}, nil
}
