package user

import (
	"encoding/json"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	logic "shop/gateway/internal/logic/user"
	"shop/gateway/internal/svc"
	"shop/gateway/internal/types"
)

func ChangeUsernameHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChangeUsernameRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// Extract user_id from context as json.Number
		userID, ok := r.Context().Value("user_id").(json.Number)
		if !ok {
			logx.Errorf("ChangeUsernameHandler: invalid user_id type, got %T", r.Context().Value("user_id"))
			httpx.Error(w, errors.New("invalid user_id in token"))
			return
		}

		// Convert json.Number to int64
		userIdInt64, err := userID.Int64()
		if err != nil {
			logx.Errorf("ChangeUsernameHandler: failed to convert user_id %v to int64: %v", userID, err)
			httpx.Error(w, errors.New("failed to convert user_id to int64"))
			return
		}
		req.UserID = userIdInt64

		l := logic.NewChangeUsernameLogic(r.Context(), svcCtx)
		resp, err := l.ChangeUsername(&req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, resp)
	}
}
