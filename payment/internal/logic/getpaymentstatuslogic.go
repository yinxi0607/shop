package logic

import (
	"context"
	"fmt"
	"shop/payment/internal/svc"
	"shop/payment/payment"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaymentStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPaymentStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaymentStatusLogic {
	return &GetPaymentStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPaymentStatusLogic) GetPaymentStatus(req *payment.GetPaymentStatusRequest) (resp *payment.GetPaymentStatusResponse, err error) {
	// 查询支付记录
	paymentQuery, err := l.svcCtx.PaymentModel.FindOneByPaymentId(l.ctx, req.PaymentId)
	if err != nil {
		return nil, fmt.Errorf("查询支付状态失败: %v", err)
	}
	if paymentQuery == nil {
		return nil, fmt.Errorf("支付记录不存在")
	}

	// 验证用户权限
	if paymentQuery.UserId != req.UserId {
		return nil, fmt.Errorf("无权访问该支付记录")
	}

	return &payment.GetPaymentStatusResponse{
		Payment: &payment.Payment{
			PaymentId:     paymentQuery.PaymentId,
			OrderId:       paymentQuery.OrderId,
			UserId:        paymentQuery.UserId,
			Amount:        paymentQuery.Amount,
			PaymentMethod: paymentQuery.PaymentMethod,
			Status:        paymentQuery.Status,
			TransactionId: paymentQuery.TransactionId.String,
			CreatedAt:     paymentQuery.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     paymentQuery.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
