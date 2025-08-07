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
	p, err := l.svcCtx.ProductModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, errors.New("product not found")
	}
	if p.DeletedAt.Valid {
		return nil, errors.New("product has been deleted")
	}

	// Update fields if provided
	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Description != nil {
		p.Description = *req.Description
	}
	if req.Detail != nil {
		//p.Detail = *req.Detail
		p.Detail = sql.NullString{
			String: *req.Detail,
			Valid:  *req.Detail != "",
		}
	}
	if req.MainImage != nil {
		p.MainImage = *req.MainImage
	}
	if req.Thumbnail != nil {
		p.Thumbnail = *req.Thumbnail
	}
	if req.Price != nil {
		p.Price = *req.Price
	}
	if req.Stock != nil {
		p.Stock = int64(*req.Stock)
	}
	if req.CategoryId != nil {
		p.CategoryId = *req.CategoryId
	}

	// Update database
	if _, err := l.svcCtx.ProductModel.Update(l.ctx, p); err != nil {
		return nil, err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("product:%d", req.Id)
	l.svcCtx.Redis.DelCtx(l.ctx, cacheKey)

	return &product.UpdateProductResponse{Success: true}, nil
}
