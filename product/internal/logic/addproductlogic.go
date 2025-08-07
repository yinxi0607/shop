package logic

import (
	"context"
	"database/sql"
	"errors"
	"shop/product/internal/svc"
	"shop/product/model"
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

func (l *AddProductLogic) AddProduct(req *product.AddProductRequest) (*product.AddProductResponse, error) {
	// Validate input
	if req.Name == "" || req.Price <= 0 || req.Stock < 0 || req.CategoryId <= 0 {
		return nil, errors.New("invalid input")
	}

	// Insert product
	p := &model.Products{
		Name:        req.Name,
		Description: req.Description,
		Detail: sql.NullString{
			String: req.Detail,
			Valid:  req.Detail != "",
		},
		MainImage:  req.MainImage,
		Thumbnail:  req.Thumbnail,
		Price:      req.Price,
		Stock:      int64(req.Stock),
		CategoryId: req.CategoryId,
		IsBanner:   req.IsBanner,
	}

	result, err := l.svcCtx.ProductModel.Insert(l.ctx, p)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &product.AddProductResponse{Id: id}, nil
}
