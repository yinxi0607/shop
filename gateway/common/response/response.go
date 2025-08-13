package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func Success(w http.ResponseWriter, data interface{}) {
	httpx.OkJsonCtx(nil, w, Body{
		Code: 0,
		Data: data,
		Msg:  "",
	})
}

func Fail(w http.ResponseWriter, code int, msg string) {
	httpx.OkJsonCtx(nil, w, Body{
		Code: code,
		Data: nil,
		Msg:  msg,
	})
}
