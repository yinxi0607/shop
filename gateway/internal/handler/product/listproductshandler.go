package product

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"shop/gateway/internal/logic/product"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
)

func ListProductsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListProductsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := product.NewListProductsLogic(r.Context(), svcCtx)
		resp, err := l.ListProducts(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}
