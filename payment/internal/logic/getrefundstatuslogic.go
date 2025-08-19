package logic

import (
	"context"
	"fmt"
	"shop/payment/internal/svc"
	"shop/payment/payment" // 使用 gRPC 生成的 payment 包

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRefundStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRefundStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRefundStatusLogic {
	return &GetRefundStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRefundStatusLogic) GetRefundStatus(req *payment.GetRefundStatusRequest) (resp *payment.GetRefundStatusResponse, err error) {
	refund, err := l.svcCtx.RefundModel.FindOneByRefundId(l.ctx, req.RefundId)
	if err != nil {
		return nil, fmt.Errorf("查询退款状态失败: %v", err)
	}
	if refund == nil {
		return nil, fmt.Errorf("退款记录不存在")
	}
	if refund.UserId != req.UserId {
		return nil, fmt.Errorf("无权访问该退款记录")
	}

	return &payment.GetRefundStatusResponse{
		Refund: &payment.Refund{
			RefundId:      refund.RefundId,
			PaymentId:     refund.PaymentId,
			OrderId:       refund.OrderId,
			UserId:        refund.UserId,
			Amount:        refund.Amount,
			Status:        refund.Status,
			TransactionId: refund.TransactionId.String,
			CreatedAt:     refund.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     refund.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
