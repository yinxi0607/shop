package cart

import (
	"context"
	"shop/cart/cart"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddToCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddToCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddToCartLogic {
	return &AddToCartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddToCartLogic) AddToCart(req *types.AddToCartRequest) (resp *types.AddToCartResponse, err error) {
	res, err := l.svcCtx.CartRpc.AddToCart(l.ctx, &cart.AddToCartRequest{
		UserId:   req.UserID,
		Pid:      req.Pid,
		Quantity: req.Quantity,
	})
	if err != nil {
		return nil, err
	}
	return &types.AddToCartResponse{Success: res.Success}, nil
}
