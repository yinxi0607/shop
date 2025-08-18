package cart

import (
	"context"
	"shop/cart/cart"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClearCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearCartLogic {
	return &ClearCartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClearCartLogic) ClearCart(req *types.ClearCartRequest) (resp *types.ClearCartResponse, err error) {
	res, err := l.svcCtx.CartRpc.ClearCart(l.ctx, &cart.ClearCartRequest{UserId: req.UserID})
	if err != nil {
		return nil, err
	}
	return &types.ClearCartResponse{Success: res.Success}, nil
}
