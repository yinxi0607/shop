package logic

import (
	"context"
	"fmt"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay/v3"
	"github.com/go-pay/xlog"
	"github.com/google/uuid"
	"shop/order/order"
	"shop/payment/internal/svc"
	"shop/payment/model"
	"shop/payment/payment"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePaymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaymentLogic {
	return &CreatePaymentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePaymentLogic) CreatePayment(req *payment.CreatePaymentRequest) (resp *payment.CreatePaymentResponse, err error) {
	// 验证订单是否存在
	orderDetail, err := l.svcCtx.OrderRpc.GetOrderDetail(l.ctx, &order.GetOrderDetailRequest{
		OrderId: req.OrderId,
		UserId:  req.UserId,
	})
	if err != nil || orderDetail.Order == nil {
		return nil, fmt.Errorf("订单不存在或无权访问")
	}

	// 生成支付ID
	paymentId := uuid.New().String()
	amount := orderDetail.Order.GetTotalPrice()

	// Redis 锁防止重复支付
	lockKey := "payment:lock:" + req.OrderId
	lock, err := l.svcCtx.Redis.SetnxExCtx(l.ctx, lockKey, "1", 10)
	if err != nil || !lock {
		return nil, fmt.Errorf("订单正在支付中")
	}
	defer l.svcCtx.Redis.DelCtx(l.ctx, lockKey)

	var paymentUrl string
	switch req.PaymentMethod {
	case "wechat":
		expire := time.Now().Add(10 * time.Minute).Format(time.RFC3339)

		bm := make(gopay.BodyMap)
		bm.Set("sp_appid", "sp_appid").
			Set("sp_mchid", "sp_mchid").
			Set("sub_mchid", "sub_mchid").
			Set("description", "测试Jsapi支付商品").
			Set("out_trade_no", paymentId).
			Set("time_expire", expire).
			Set("notify_url", "https://www.fmm.ink").
			SetBodyMap("amount", func(bm gopay.BodyMap) {
				bm.Set("total", int64(amount*100)).
					Set("currency", "CNY")
			}).
			SetBodyMap("payer", func(bm gopay.BodyMap) {
				bm.Set("sp_openid", "asdas")
			})
		wxRsp, err := l.svcCtx.WechatClient.V3TransactionNative(l.ctx, bm)
		if err != nil {
			l.Logger.Error(err)
			return
		}
		xlog.Errorf("wxRsp:%s", wxRsp.Error)
		paymentUrl = wxRsp.Response.CodeUrl
	case "alipay":

		bm := make(gopay.BodyMap)
		bm.Set("subject", "pre create order").
			Set("out_trade_no", paymentId).
			Set("total_amount", fmt.Sprintf("%.2f", amount)).
			Set(alipay.HeaderAppAuthToken, l.svcCtx.AlipayClient.AppAuthToken) // 如果需要，可以设置自定义应用授权

		// 创建订单
		aliRsp, err := l.svcCtx.AlipayClient.TradePrecreate(l.ctx, bm)
		if err != nil {
			xlog.Errorf("client.TradePrecreate(), err:%v", err)
			return
		}
		l.Logger.Infof("aliRsp:%s", aliRsp)

		l.Logger.Infof("aliRsp.QrCode:", aliRsp.QrCode)
		paymentUrl = aliRsp.QrCode

	default:
		return nil, fmt.Errorf("不支持的支付方式")
	}

	// 保存支付记录
	_, err = l.svcCtx.PaymentModel.Insert(l.ctx, &model.Payments{
		PaymentId:     paymentId,
		OrderId:       req.OrderId,
		UserId:        req.UserId,
		Amount:        amount,
		PaymentMethod: req.PaymentMethod,
		Status:        "pending",
	})
	if err != nil {
		return nil, err
	}

	// 缓存支付状态
	err = l.svcCtx.Redis.SetexCtx(l.ctx, "payment:"+req.OrderId, "pending", 3600)
	if err != nil {
		return nil, err
	}

	return &payment.CreatePaymentResponse{
		PaymentId:  paymentId,
		PaymentUrl: paymentUrl,
	}, nil
}
