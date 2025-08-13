package order

import (
	"net/http"
	"shop/gateway/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
	"shop/gateway/internal/logic/order"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
)

func SeckillOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SeckillOrderRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Fail(w, 1000, err.Error())
			return
		}

		l := order.NewSeckillOrderLogic(r.Context(), svcCtx)
		resp, err := l.SeckillOrder(&req)
		if err != nil {
			//httpx.ErrorCtx(r.Context(), w, err)
			response.Fail(w, 1000, err.Error())
		} else {
			//httpx.OkJsonCtx(r.Context(), w, resp)
			response.Success(w, resp)
		}
	}
}
