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

func (l *ListProductsLogic) ListProducts(req *product.ListProductsRequest) (*product.ListProductsResponse, error) {
	// Check Redis cache
	cacheKey := fmt.Sprintf("products:list:%d:%d:%v:%v:%v", req.Page, req.PageSize, req.CategoryId, req.MinPrice, req.MaxPrice)
	cached, err := l.svcCtx.Redis.GetCtx(l.ctx, cacheKey)
	if err == nil {
		var resp product.ListProductsResponse
		if err := json.Unmarshal([]byte(cached), &resp); err == nil {
			return &resp, nil
		}
	}

	// Query database
	var categoryId *int64
	var minPrice, maxPrice *float64
	if req.CategoryId != nil {
		categoryId = req.CategoryId
	}
	if req.MinPrice != nil {
		minPrice = req.MinPrice
	}
	if req.MaxPrice != nil {
		maxPrice = req.MaxPrice
	}

	products, total, err := l.svcCtx.ProductModel.List(l.ctx, req.Page, req.PageSize, categoryId, minPrice, maxPrice)
	if err != nil {
		return nil, err
	}

	resp := &product.ListProductsResponse{Total: total}
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
		l.svcCtx.Redis.SetexCtx(l.ctx, cacheKey, string(data), 5*60) // Cache for 5 minutes
	}

	return resp, nil
}
