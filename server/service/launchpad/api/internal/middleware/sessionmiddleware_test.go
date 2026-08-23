package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"o1-launchpad/common/session"
)

func TestSessionMiddlewareRequiresValidHttpOnlySessionValue(t *testing.T) {
	secret := "test-secret"
	expiresAt := time.Now().Add(time.Hour)
	token, err := session.Encode(secret, session.Claims{Address: "0xabc", ExpiresAt: expiresAt.Unix()})
	if err != nil {
		t.Fatal(err)
	}

	handler := SessionMiddleware(secret)(func(w http.ResponseWriter, r *http.Request) {
		address, ok := AddressFromContext(r.Context())
		if !ok || address != "0xabc" {
			t.Fatalf("session address = %q, %v", address, ok)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	missing := httptest.NewRecorder()
	handler(missing, httptest.NewRequest(http.MethodPost, "/v1/launchpads", nil))
	if missing.Code != http.StatusUnauthorized {
		t.Fatalf("missing session status = %d", missing.Code)
	}

	request := httptest.NewRequest(http.MethodPost, "/v1/launchpads", nil)
	request.AddCookie(&http.Cookie{Name: "o1_session", Value: token})
	response := httptest.NewRecorder()
	handler(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("valid session status = %d", response.Code)
	}
}
