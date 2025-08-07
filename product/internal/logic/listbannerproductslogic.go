package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"shop/product/internal/svc"
	"shop/product/product"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListBannerProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListBannerProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBannerProductsLogic {
	return &ListBannerProductsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListBannerProductsLogic) ListBannerProducts(req *product.ListBannerProductsRequest) (*product.ListBannerProductsResponse, error) {
	// Check Redis cache
	cacheKey := fmt.Sprintf("products:banner:%d", req.Limit)
	cached, err := l.svcCtx.Redis.GetCtx(l.ctx, cacheKey)
	if err == nil {
		var resp product.ListBannerProductsResponse
		if err := json.Unmarshal([]byte(cached), &resp); err == nil {
			return &resp, nil
		}
	}

	// Query database
	products, err := l.svcCtx.ProductModel.ListBanner(l.ctx, req.Limit)
	if err != nil {
		return nil, err
	}

	resp := &product.ListBannerProductsResponse{}
	for _, p := range products {
		resp.Products = append(resp.Products, &product.Product{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Detail:      p.Detail.String,
			MainImage:   p.MainImage,
			Thumbnail:   p.Thumbnail,
			Price:       p.Price,
			Stock:       int32(p.Stock),
			CategoryId:  p.CategoryId,
			CreatedAt:   p.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
		})
	}

	// Cache to Redis
	if data, err := json.Marshal(resp); err == nil {
		l.svcCtx.Redis.SetexCtx(l.ctx, cacheKey, string(data), 5*60)
	}

	return resp, nil
}
