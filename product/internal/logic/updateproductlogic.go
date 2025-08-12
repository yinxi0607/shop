package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"shop/product/internal/svc"
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

func (l *UpdateProductLogic) UpdateProduct(req *product.UpdateProductRequest) (*product.UpdateProductResponse, error) {
	// Find existing product
	p, err := l.svcCtx.ProductModel.FindOne(l.ctx, req.Pid)
	if err != nil {
		return nil, errors.New("product not found")
	}
	if p.DeletedAt.Valid {
		return nil, errors.New("product has been deleted")
	}

	// Update fields if provided
	if req.Name != nil && *req.Name != "" {
		p.Name = *req.Name
	}
	if req.Description != nil && *req.Description != "" {
		p.Description = *req.Description
	}
	if req.Detail != nil && *req.Detail != "" {
		//p.Detail = *req.Detail
		p.Detail = sql.NullString{
			String: *req.Detail,
			Valid:  *req.Detail != "",
		}
	}
	if req.MainImage != nil && *req.MainImage != "" {
		p.MainImage = *req.MainImage
	}
	if req.Thumbnail != nil && *req.Thumbnail != "" {
		p.Thumbnail = *req.Thumbnail
	}
	if req.Price != nil && *req.Price != 0 {
		p.Price = *req.Price
	}
	if req.Stock != nil && *req.Stock != 0 {
		p.Stock = int64(*req.Stock)
	}
	if req.CategoryId != nil && *req.CategoryId != "" {
		p.CategoryId = *req.CategoryId
	}

	// Update database
	if err = l.svcCtx.ProductModel.Update(l.ctx, p); err != nil {
		return nil, err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("product:%s", req.Pid)
	_, err = l.svcCtx.Redis.DelCtx(l.ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	return &product.UpdateProductResponse{Success: true}, nil
}
