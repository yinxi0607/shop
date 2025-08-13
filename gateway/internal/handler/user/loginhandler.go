package user

import (
	"net/http"
	"shop/gateway/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
	"shop/gateway/internal/logic/user"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Fail(w, 1000, err.Error())
			return
		}

		l := user.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		if err != nil {
			response.Fail(w, 1000, err.Error())
		} else {
			response.Success(w, resp)
		}
	}
}
