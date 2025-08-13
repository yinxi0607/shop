package user

import (
	"net/http"
	"shop/gateway/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
	"shop/gateway/internal/logic/user"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
)

func ChangeRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChangeRoleRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Fail(w, 10000, err.Error())
			return
		}

		l := user.NewChangeRoleLogic(r.Context(), svcCtx)
		resp, err := l.ChangeRole(&req)
		if err != nil {
			response.Fail(w, 10000, err.Error())
		} else {
			response.Success(w, resp)
		}
	}
}
