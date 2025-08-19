package logic

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-pay/gopay"
	"github.com/go-pay/util"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"shop/order/order"
	"shop/payment/internal/svc"
	"shop/payment/model"
	"shop/payment/payment"

	"github.com/IBM/sarama"
	"github.com/zeromicro/go-zero/core/logx"
)

type RefundLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundLogic {
	return &RefundLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefundLogic) Refund(req *payment.RefundRequest) (resp *payment.RefundResponse, err error) {
	// 验证支付记录
	paymentQuery, err := l.svcCtx.PaymentModel.FindOneByPaymentId(l.ctx, req.PaymentId)
	if err != nil || paymentQuery == nil {
		return nil, fmt.Errorf("支付记录不存在")
	}
	if paymentQuery.UserId != req.UserId {
		return nil, fmt.Errorf("无权发起退款")
	}
	if paymentQuery.Status != "completed" {
		return nil, fmt.Errorf("支付状态不支持退款: %s", paymentQuery.Status)
	}
	if req.Amount > paymentQuery.Amount {
		return nil, fmt.Errorf("退款金额超限")
	}

	// Redis 锁防止重复退款
	lockKey := "refund:lock:" + req.PaymentId
	lock, err := l.svcCtx.Redis.SetnxExCtx(l.ctx, lockKey, "1", 10)
	if err != nil || !lock {
		return nil, fmt.Errorf("退款处理中")
	}
	defer func(Redis *redis.Redis, ctx context.Context, keys string) {
		_, err = Redis.DelCtx(ctx, keys)
		if err != nil {

		}
	}(l.svcCtx.Redis, l.ctx, lockKey)

	// 生成退款ID
	refundId := uuid.New().String()

	// 验证订单状态
	orderResp, err := l.svcCtx.OrderRpc.GetOrderDetail(l.ctx, &order.GetOrderDetailRequest{
		OrderId: paymentQuery.OrderId,
		UserId:  req.UserId,
	})
	if err != nil || orderResp.Order == nil {
		return nil, fmt.Errorf("订单不存在或无权访问")
	}
	if orderResp.Order.Status != "paid" {
		return nil, fmt.Errorf("订单状态不支持退款: %s", orderResp.Order.Status)
	}

	var transactionId string
	switch paymentQuery.PaymentMethod {
	case "wechat":
		//生成订单流水号
		s := util.RandomString(64)
		orderNo := fmt.Sprintf("CX-%s", s)

		// 初始化 BodyMap
		bm := make(gopay.BodyMap)
		// 选填 商户订单号（支付后返回的，一般是以42000开头）
		bm.Set("transaction_id", paymentQuery.TransactionId).
			Set("sign_type", "MD5").
			// 必填 退款订单号（程序员定义的）
			Set("out_refund_no", orderNo).
			// 选填 退款描述
			Set("reason", "这是一退款操作").
			//Set("notify_url", l.svcCtx.WechatClient.notif).
			SetBodyMap("amount", func(bm gopay.BodyMap) {
				// 退款金额:单位是分
				bm.Set("refund", paymentQuery.Amount). //实际退款金额
									Set("total", paymentQuery.Amount). // 折扣前总金额（不是实际退款数）
									Set("currency", "CNY")
			})
		//请求申请退款（沙箱环境下，证书路径参数可传空）
		//    body：参数Body
		refund, err := l.svcCtx.WechatClient.V3Refund(l.ctx, bm)
		if err != nil {
			l.Logger.Error(err)
			return
		}

		// 将非正常退款异常记录
		// 返回：404 > {"code":"RESOURCE_NOT_EXISTS","message":"订单不存在"}
		if refund.Code == 404 || refund.Code == 400 || refund.Code == 403 {
			// 这里时对非正常退款的一些处理message，我们将code统一使用自定义的，然后把message抛出去

			l.Logger.Infof(fmt.Sprintf("code:%d,message:%s", 50000, refund.ErrResponse.Message))
			return
		}
		transactionId = refund.Response.TransactionId

	case "alipay":
		//配置公共参

		//请求参数
		bm := make(gopay.BodyMap)
		bm.Set("out_trade_no", "GZ201907301420334577")
		bm.Set("refund_amount", "5")
		bm.Set("refund_reason", "测试退款")
		//发起退款请求
		aliRsp, err := l.svcCtx.AlipayClient.TradeRefund(l.ctx, bm)
		if err != nil {
			l.Logger.Error("err:", err)
			return
		}
		l.Logger.Infof("aliRsp:", *aliRsp)
		transactionId = aliRsp.TradeNo
	default:
		return nil, fmt.Errorf("不支持的支付方式")
	}

	// 保存退款记录
	_, err = l.svcCtx.RefundModel.Insert(l.ctx, &model.Refunds{
		RefundId:      refundId,
		PaymentId:     req.PaymentId,
		OrderId:       paymentQuery.OrderId,
		UserId:        req.UserId,
		Amount:        req.Amount,
		Status:        "pending",
		TransactionId: sql.NullString{String: transactionId, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("保存退款记录失败: %v", err)
	}

	// 缓存退款状态
	err = l.svcCtx.Redis.SetexCtx(l.ctx, "refund:"+req.PaymentId, "pending", 3600)
	if err != nil {
		return nil, fmt.Errorf("缓存退款状态失败: %v", err)
	}

	// 发送 Kafka 通知
	msg := &struct {
		Type          string  `json:"type"`
		PaymentId     string  `json:"payment_id"`
		RefundId      string  `json:"refund_id"`
		OrderId       string  `json:"order_id"`
		Status        string  `json:"status"`
		Amount        float64 `json:"amount"`
		TransactionId string  `json:"transaction_id"`
	}{
		Type:          "refund",
		PaymentId:     req.PaymentId,
		RefundId:      refundId,
		OrderId:       paymentQuery.OrderId,
		Status:        "pending",
		Amount:        req.Amount,
		TransactionId: transactionId,
	}
	_, _, err = l.svcCtx.KafkaProducer.SendMessage(&sarama.ProducerMessage{
		Topic: l.svcCtx.Config.Kafka.Topic,
		Value: sarama.StringEncoder(fmt.Sprintf(`{"type":"%s","payment_id":"%s","refund_id":"%s","order_id":"%s","status":"%s","amount":%.2f,"transaction_id":"%s"}`,
			msg.Type, msg.PaymentId, msg.RefundId, msg.OrderId, msg.Status, msg.Amount, msg.TransactionId)),
	})
	if err != nil {
		return nil, fmt.Errorf("发送退款通知失败: %v", err)
	}

	return &payment.RefundResponse{RefundId: refundId}, nil
}
