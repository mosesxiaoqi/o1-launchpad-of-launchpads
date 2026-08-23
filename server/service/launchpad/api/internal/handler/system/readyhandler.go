package system

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"o1-launchpad/service/launchpad/api/internal/logic/system"
	"o1-launchpad/service/launchpad/api/internal/svc"
)

func ReadyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := system.NewReadyLogic(r.Context(), svcCtx)
		resp, err := l.Ready()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
