package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"msg/api/internal/logic"
	"msg/api/internal/svc"
	"msg/api/internal/types"
)

func TransHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TransRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewTransLogic(r.Context(), svcCtx)
		resp, err := l.Trans(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
