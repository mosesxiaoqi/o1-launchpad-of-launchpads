package middleware

import (
	"context"
	"net/http"
	"time"

	"o1-launchpad/common/session"
)

type addressKey struct{}

func SessionMiddleware(secret string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("o1_session")
			if err != nil {
				http.Error(w, "authentication required", http.StatusUnauthorized)
				return
			}
			claims, err := session.Decode(secret, cookie.Value, time.Now())
			if err != nil {
				http.Error(w, "invalid session", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), addressKey{}, claims.Address)
			next(w, r.WithContext(ctx))
		}
	}
}

func AddressFromContext(ctx context.Context) (string, bool) {
	address, ok := ctx.Value(addressKey{}).(string)
	return address, ok
}
