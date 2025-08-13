package product

import (
	"net/http"
	"shop/gateway/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
	"shop/gateway/internal/logic/product"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
)

func ListRecommendedProductsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListRecommendedProductsRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Fail(w, 1000, err.Error())
			return
		}

		l := product.NewListRecommendedProductsLogic(r.Context(), svcCtx)
		resp, err := l.ListRecommendedProducts(&req)
		if err != nil {
			response.Fail(w, 1000, err.Error())
		} else {
			response.Success(w, resp)
		}
	}
}
