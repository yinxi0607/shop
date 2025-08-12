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

type GetProductDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductDetailLogic {
	return &GetProductDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductDetailLogic) GetProductDetail(req *product.GetProductDetailRequest) (*product.GetProductDetailResponse, error) {
	// Check Redis cache
	cacheKey := fmt.Sprintf("product:%s", req.Pid)
	cached, err := l.svcCtx.Redis.GetCtx(l.ctx, cacheKey)
	if err == nil {
		var prod product.Product
		if err := json.Unmarshal([]byte(cached), &prod); err == nil {
			return &product.GetProductDetailResponse{Product: &prod}, nil
		}
	}

	// Query database
	p, err := l.svcCtx.ProductModel.FindOne(l.ctx, req.Pid)
	if err != nil {
		return nil, err
	}

	resp := &product.GetProductDetailResponse{
		Product: &product.Product{
			Pid:         p.Pid,
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
		},
	}

	// Cache to Redis
	if data, err := json.Marshal(resp.Product); err == nil {
		l.svcCtx.Redis.SetexCtx(l.ctx, cacheKey, string(data), 5*60) // Cache for 5 minutes
	}

	return resp, nil
}
