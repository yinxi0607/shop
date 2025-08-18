package cart

import (
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"shop/gateway/common/response"
	"shop/gateway/internal/logic/cart"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ClearCartHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ClearCartRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Fail(w, 1000, err.Error())
			return
		}
		userID, ok := r.Context().Value("user_id").(string)
		if !ok {
			logx.Errorf("ClearCartHandler: invalid user_id type, got %T", r.Context().Value("user_id"))
			response.Fail(w, 10000, "invalid user_id")
			return
		}
		req.UserID = userID
		l := cart.NewClearCartLogic(r.Context(), svcCtx)
		resp, err := l.ClearCart(&req)
		if err != nil {
			response.Fail(w, 1000, err.Error())
		} else {
			response.Success(w, resp)
		}
	}
}
