package user

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"shop/gateway/common/response"
	"shop/gateway/internal/logic/user"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
)

func GetUserInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserInfoRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Fail(w, 10000, err.Error())
			return
		}

		// 从 JWT 获取 user_id
		userIDValue := r.Context().Value("user_id")
		logx.Infof("GetUserInfoHandler: user_id=%v, type=%T", userIDValue, userIDValue)

		userID, ok := r.Context().Value("user_id").(string)
		if !ok {
			logx.Errorf("GetUserInfoHandler: invalid user_id type, got %T", r.Context().Value("user_id"))
			response.Fail(w, 10000, "invalid user_id")
			return
		}

		req.UserID = userID

		l := user.NewGetUserInfoLogic(r.Context(), svcCtx)
		resp, err := l.GetUserInfo(&req)
		if err != nil {
			response.Fail(w, 10000, err.Error())
			return
		}
		response.Success(w, resp)
	}
}
