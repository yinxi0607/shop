package product

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
	"shop/product/product"
)

type ListRecommendedProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRecommendedProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRecommendedProductsLogic {
	return &ListRecommendedProductsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRecommendedProductsLogic) ListRecommendedProducts(req *types.ListRecommendedProductsRequest) (resp *types.ListRecommendedProductsResponse, err error) {
	// 检查 Redis 缓存
	cacheKey := fmt.Sprintf("recommended_products:limit:%d", req.Limit)
	cached, err := l.svcCtx.Redis.GetCtx(l.ctx, cacheKey)
	if err == nil && cached != "" {
		var products []types.Product
		if err := json.Unmarshal([]byte(cached), &products); err == nil {
			logx.Infof("ListRecommendedProductsLogic: cache hit for key %s", cacheKey)
			return &types.ListRecommendedProductsResponse{
				Products: products,
			}, nil
		}
	}

	// 调用商品服务的 gRPC 接口
	productResp, err := l.svcCtx.ProductRpc.ListRecommendedProducts(l.ctx, &product.ListRecommendedProductsRequest{
		Limit: req.Limit,
	})
	if err != nil {
		logx.Errorf("ListRecommendedProductsLogic: failed to call ProductRpc.ListRecommendedProducts: %v", err)
		return nil, err
	}

	// 转换 gRPC 响应
	products := make([]types.Product, len(productResp.Products))
	for i, p := range productResp.Products {
		products[i] = types.Product{
			Pid:         p.Pid,
			Name:        p.Name,
			Description: p.Description,
			Detail:      p.Detail,
			MainImage:   p.MainImage,
			Thumbnail:   p.Thumbnail,
			Price:       p.Price,
			Stock:       p.Stock,
			CategoryID:  p.CategoryId,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		}
	}

	// 缓存结果（TTL 5 分钟）
	cacheData, _ := json.Marshal(products)
	err = l.svcCtx.Redis.SetexCtx(l.ctx, cacheKey, string(cacheData), 300)
	if err != nil {
		logx.Errorf("ListRecommendedProductsLogic: failed to set cache %s: %v", cacheKey, err)
		// 不返回错误，仅记录日志
	}

	return &types.ListRecommendedProductsResponse{
		Products: products,
	}, nil
}
