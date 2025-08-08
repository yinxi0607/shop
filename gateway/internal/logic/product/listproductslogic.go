package product

import (
	"context"
	"encoding/json"
	"fmt"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
	"shop/product/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductsLogic {
	return &ListProductsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListProductsLogic) ListProducts(req *types.ListProductsRequest) (resp *types.ListProductsResponse, err error) {
	// 检查 Redis 缓存
	cacheKey := fmt.Sprintf("products:page:%d:size:%d:cat:%d:min:%f:max:%f", req.Page, req.PageSize, req.CategoryID, req.MinPrice, req.MaxPrice)
	cached, err := l.svcCtx.Redis.Get(cacheKey)
	if err == nil && cached != "" {
		var products []types.Product
		if err := json.Unmarshal([]byte(cached), &products); err == nil {
			logx.Infof("ListProductsLogic: cache hit for key %s", cacheKey)
			return &types.ListProductsResponse{
				Products: products,
				Total:    int32(len(products)), // 假设缓存包含总数
			}, nil
		}
	}
	// 调用商品服务的 gRPC 接口
	productResp, err := l.svcCtx.ProductRpc.ListProducts(l.ctx, &product.ListProductsRequest{
		Page:       req.Page,
		PageSize:   req.PageSize,
		CategoryId: req.CategoryID,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
	})
	if err != nil {
		logx.Errorf("ListProductsLogic: failed to call ProductRpc.ListProducts: %v", err)
		return nil, err
	}

	// 转换 gRPC 响应
	products := make([]types.Product, len(productResp.Products))
	for i, p := range productResp.Products {
		products[i] = types.Product{
			ID:          p.Id,
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
	err = l.svcCtx.Redis.Setex(cacheKey, string(cacheData), 300)
	if err != nil {
		logx.Errorf("ListProductsLogic: failed to set cache %s: %v", cacheKey, err)
		// 不返回错误，仅记录日志
	}

	return &types.ListProductsResponse{
		Products: products,
		Total:    productResp.Total,
	}, nil
}
