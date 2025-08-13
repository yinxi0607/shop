package order

import (
	"net/http"
	"shop/gateway/common/response"
	"shop/gateway/internal/logic/order"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateOrderStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateOrderStatusRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Fail(w, 1000, err.Error())
			return
		}

		l := order.NewUpdateOrderStatusLogic(r.Context(), svcCtx)
		resp, err := l.UpdateOrderStatus(&req)
		if err != nil {
			response.Fail(w, 1000, err.Error())
		} else {
			response.Success(w, resp)
		}
	}
}
