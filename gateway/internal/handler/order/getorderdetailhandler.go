package order

import (
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"shop/gateway/common/response"
	"shop/gateway/internal/logic/order"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetOrderDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetOrderDetailRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Fail(w, 1000, err.Error())
			return
		}
		// 从 JWT 获取 user_id
		userID, ok := r.Context().Value("user_id").(string)
		if !ok {
			logx.Errorf("AddProductHandler: invalid user_id type, got %T", r.Context().Value("user_id"))
			response.Fail(w, 10000, "invalid user_id")
			return
		}
		req.UserID = userID
		l := order.NewGetOrderDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetOrderDetail(&req)
		if err != nil {
			response.Fail(w, 1000, err.Error())
		} else {
			response.Success(w, resp)
		}
	}
}
