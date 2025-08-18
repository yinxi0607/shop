package cart

import (
	"context"
	"shop/cart/cart"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCartLogic {
	return &GetCartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCartLogic) GetCart(req *types.GetCartRequest) (resp *types.GetCartResponse, err error) {
	res, err := l.svcCtx.CartRpc.GetCart(l.ctx, &cart.GetCartRequest{UserId: req.UserID})
	if err != nil {
		return nil, err
	}
	items := make([]types.CartItem, len(res.Items))
	for i, item := range res.Items {
		items[i] = types.CartItem{Pid: item.Pid, Quantity: item.Quantity}
	}
	return &types.GetCartResponse{
		CartID: res.CartId,
		Items:  items,
	}, nil
}
