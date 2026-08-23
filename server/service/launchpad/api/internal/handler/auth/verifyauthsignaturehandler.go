package auth

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"o1-launchpad/service/launchpad/api/internal/logic/auth"
	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"
)

func VerifyAuthSignatureHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VerifyAuthSignatureRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := auth.NewVerifyAuthSignatureLogic(r.Context(), svcCtx)
		resp, err := l.VerifyAuthSignature(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
