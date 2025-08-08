package product

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"shop/gateway/internal/logic/product"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
)

func GetProductDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetProductDetailRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := product.NewGetProductDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetProductDetail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
