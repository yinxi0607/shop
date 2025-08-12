package product

import (
	"encoding/json"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"shop/gateway/internal/logic/product"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
)

func UpdateProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateProductRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 从 JWT 获取 user_id
		userID, ok := r.Context().Value("user_id").(json.Number)
		if !ok {
			logx.Errorf("UpdateProductHandler: invalid user_id type, got %T", r.Context().Value("user_id"))
			httpx.Error(w, errors.New("invalid user_id in token"))
			return
		}
		userIdInt64, err := userID.Int64()
		if err != nil {
			logx.Errorf("UpdateProductHandler: failed to convert user_id %v to int64: %v", userID, err)
			httpx.Error(w, errors.New("failed to convert user_id to int64"))
			return
		}

		// 管理员权限检查
		if !isAdmin(userIdInt64, svcCtx.Config.AdminUser) {
			logx.Errorf("UpdateProductHandler: user_id %d is not admin", userIdInt64)
			httpx.Error(w, errors.New("admin access required"))
			return
		}

		l := product.NewUpdateProductLogic(r.Context(), svcCtx)
		resp, err := l.UpdateProduct(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}
