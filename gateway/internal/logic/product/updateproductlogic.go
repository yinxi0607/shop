package product

import (
	"context"
	"fmt"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
	"shop/product/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductLogic {
	return &UpdateProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProductLogic) UpdateProduct(req *types.UpdateProductRequest) (resp *types.UpdateProductResponse, err error) {
	// 调用商品服务的 gRPC 接口
	_, err = l.svcCtx.ProductRpc.UpdateProduct(l.ctx, &product.UpdateProductRequest{
		Pid:         req.Pid,
		Name:        &req.Name,
		Description: &req.Description,
		Detail:      &req.Detail,
		MainImage:   &req.MainImage,
		Thumbnail:   &req.Thumbnail,
		Price:       &req.Price,
		Stock:       &req.Stock,
		CategoryId:  &req.CategoryID,
	})
	if err != nil {
		logx.Errorf("UpdateProductLogic: failed to call ProductRpc.UpdateProduct: %v", err)
		return nil, err
	}

	// 失效 Redis 缓存
	cacheKey := fmt.Sprintf("product:%s", req.Pid)
	_, err = l.svcCtx.Redis.Del(cacheKey)
	if err != nil {
		logx.Errorf("UpdateProductLogic: failed to delete cache %s: %v", cacheKey, err)
		// 不返回错误，仅记录日志
	}

	return &types.UpdateProductResponse{
		Success: true,
	}, nil
}
