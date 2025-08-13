package product

import (
	"encoding/json"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"shop/gateway/common/response"
	"shop/gateway/internal/logic/product"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
)

func AddProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddProductRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Fail(w, 10000, err.Error())
			return
		}

		// 从 JWT 获取 user_id
		userID, ok := r.Context().Value("user_id").(json.Number)
		if !ok {
			logx.Errorf("AddProductHandler: invalid user_id type, got %T", r.Context().Value("user_id"))
			response.Fail(w, 10000, "invalid user_id")
			return
		}
		userIdInt64, err := userID.Int64()
		if err != nil {
			logx.Errorf("AddProductHandler: failed to convert user_id %v to int64: %v", userID, err)
			httpx.Error(w, errors.New("failed to convert user_id to int64"))
			return
		}

		role, ok := r.Context().Value("role").(string)
		if !ok {
			logx.Errorf("AddProductHandler: invalid user_id type, got %T", r.Context().Value("user_id"))
			httpx.Error(w, errors.New("invalid user_id in token"))
			return
		}

		// 管理员权限检查
		if !isAdmin(role) {
			logx.Errorf("AddProductHandler: user_id %d,role %s is not admin", userIdInt64)
			httpx.Error(w, errors.New("admin access required"))
			return
		}

		l := product.NewAddProductLogic(r.Context(), svcCtx)
		resp, err := l.AddProduct(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}

func isAdmin(role string) bool {
	return role == "admin" // Replace with actual role check (e.g., query users.role)
}
