package product

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"shop/gateway/common/response"
	"shop/gateway/common/utils"
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
		userID, ok := r.Context().Value("user_id").(string)
		if !ok {
			logx.Errorf("AddProductHandler: invalid user_id type, got %T", r.Context().Value("user_id"))
			response.Fail(w, 10000, "invalid user_id")
			return
		}
		role, ok := r.Context().Value("role").(string)
		if !ok {
			logx.Errorf("AddProductHandler: invalid user_id type, got %T", r.Context().Value("user_id"))
			response.Fail(w, 10000, "invalid user_id")
			return
		}

		// 管理员权限检查
		if !utils.IsAdmin(role) {
			logx.Errorf("AddProductHandler: user_id %d,role %s is not admin", userID)
			response.Fail(w, 10000, "user_id is not admin")
			return
		}

		l := product.NewAddProductLogic(r.Context(), svcCtx)
		resp, err := l.AddProduct(&req)
		if err != nil {
			response.Fail(w, 10000, err.Error())
			return
		}

		response.Success(w, resp)
	}
}
