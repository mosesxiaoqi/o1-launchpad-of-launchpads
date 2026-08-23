package auth

import (
	"net/http"
	"time"

	"o1-launchpad/common/session"
	"o1-launchpad/service/launchpad/api/internal/logic/auth"
	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
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
			expiresAt, err := time.Parse(time.RFC3339, resp.ExpiresAt)
			if err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}
			token, err := session.Encode(svcCtx.Config.Auth.SessionSecret, session.Claims{
				Address: resp.Address, ExpiresAt: expiresAt.Unix(),
			})
			if err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name: "o1_session", Value: token, Path: "/", Expires: expiresAt, MaxAge: int(time.Until(expiresAt).Seconds()),
				HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode,
			})
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
