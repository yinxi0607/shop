package cart

import (
	"context"
	"shop/cart/cart"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCartLogic {
	return &UpdateCartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCartLogic) UpdateCart(req *types.UpdateCartRequest) (resp *types.UpdateCartResponse, err error) {
	res, err := l.svcCtx.CartRpc.UpdateCart(l.ctx, &cart.UpdateCartRequest{
		UserId:   req.UserID,
		Pid:      req.Pid,
		Quantity: req.Quantity,
	})
	if err != nil {
		return nil, err
	}
	return &types.UpdateCartResponse{Success: res.Success}, nil
}
