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

func (l *GetProductDetailLogic) GetProductDetail(req *types.GetProductDetailRequest) (resp *types.GetProductDetailResponse, err error) {
	// 检查 Redis 缓存
	cacheKey := fmt.Sprintf("product:%d", req.ID)
	cached, err := l.svcCtx.Redis.GetCtx(l.ctx, cacheKey)
	if err == nil && cached != "" {
		var p types.Product
		if err := json.Unmarshal([]byte(cached), &p); err == nil {
			logx.Infof("GetProductDetailLogic: cache hit for key %s", cacheKey)
			return &types.GetProductDetailResponse{
				Product: p,
			}, nil
		}
	}

	// 调用商品服务的 gRPC 接口
	productResp, err := l.svcCtx.ProductRpc.GetProductDetail(l.ctx, &product.GetProductDetailRequest{
		Id: req.ID,
	})
	if err != nil {
		logx.Errorf("GetProductDetailLogic: failed to call ProductRpc.GetProductDetail: %v", err)
		return nil, err
	}

	// 转换 gRPC 响应
	p := types.Product{
		ID:          productResp.Product.Id,
		Name:        productResp.Product.Name,
		Description: productResp.Product.Description,
		Detail:      productResp.Product.Detail,
		MainImage:   productResp.Product.MainImage,
		Thumbnail:   productResp.Product.Thumbnail,
		Price:       productResp.Product.Price,
		Stock:       productResp.Product.Stock,
		CategoryID:  productResp.Product.CategoryId,
		CreatedAt:   productResp.Product.CreatedAt,
		UpdatedAt:   productResp.Product.UpdatedAt,
	}

	// 缓存结果（TTL 5 分钟）
	cacheData, _ := json.Marshal(p)
	err = l.svcCtx.Redis.SetexCtx(l.ctx, cacheKey, string(cacheData), 300)
	if err != nil {
		logx.Errorf("GetProductDetailLogic: failed to set cache %s: %v", cacheKey, err)
		// 不返回错误，仅记录日志
	}

	return &types.GetProductDetailResponse{
		Product: p,
	}, nil
}
