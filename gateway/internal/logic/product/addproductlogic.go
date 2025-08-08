package product

import (
	"context"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
	"shop/product/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddProductLogic {
	return &AddProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddProductLogic) AddProduct(req *types.AddProductRequest) (resp *types.AddProductResponse, err error) {
	// 调用商品服务的 gRPC 接口
	productResp, err := l.svcCtx.ProductRpc.AddProduct(l.ctx, &product.AddProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Detail:      req.Detail,
		MainImage:   req.MainImage,
		Thumbnail:   req.Thumbnail,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryId:  req.CategoryID,
		IsBanner:    req.IsBanner,
	})
	if err != nil {
		logx.Errorf("AddProductLogic: failed to call ProductRpc.AddProduct: %v", err)
		return nil, err
	}

	return &types.AddProductResponse{
		ID: productResp.Id,
	}, nil
}
