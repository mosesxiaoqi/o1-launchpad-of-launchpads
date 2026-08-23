package fee

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"o1-launchpad/service/launchpad/api/internal/logic/fee"
	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"
)

func GetFeesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetFeesRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := fee.NewGetFeesLogic(r.Context(), svcCtx)
		resp, err := l.GetFees(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
