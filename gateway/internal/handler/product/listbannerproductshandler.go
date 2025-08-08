package product

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"shop/gateway/internal/logic/product"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
)

func ListBannerProductsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListBannerProductsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := product.NewListBannerProductsLogic(r.Context(), svcCtx)
		resp, err := l.ListBannerProducts(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
