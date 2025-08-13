package user

import (
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"shop/gateway/common/response"
	"shop/gateway/common/utils"

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

		// 从 JWT 获取 user_id
		userID, ok := r.Context().Value("user_id").(string)
		if !ok {
			logx.Errorf("ChangePasswordHandler: invalid user_id type, got %T", r.Context().Value("user_id"))
			response.Fail(w, 10000, "invalid user_id")
			return
		}

		req.UserID = userID
		role, ok := r.Context().Value("role").(string)
		if !ok {
			logx.Errorf("AddProductHandler: invalid user_id type, got %T", r.Context().Value("user_id"))
			response.Fail(w, 1000, "invalid user_id")
			return
		}

		// 管理员权限检查
		if !utils.IsAdmin(role) {
			logx.Errorf("Update user Role: user_id %d is not admin", userID)
			response.Fail(w, 1000, "user_id is not admin")
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
