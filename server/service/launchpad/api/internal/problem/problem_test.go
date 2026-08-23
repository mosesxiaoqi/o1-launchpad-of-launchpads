package problem

import (
	"errors"
	"net/http"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestFromGRPC(t *testing.T) {
	code, body := From(status.Error(codes.InvalidArgument, "invalid cursor"))
	if code != http.StatusBadRequest || body.Status != http.StatusBadRequest || body.Code != "invalid_argument" {
		t.Fatalf("From() = %d, %#v", code, body)
	}
	if body.Detail != "invalid cursor" || body.Type != "about:blank" {
		t.Fatalf("unexpected problem body: %#v", body)
	}
}

func TestFromLocalBindingError(t *testing.T) {
	code, body := From(errors.New("missing currency"))
	if code != http.StatusBadRequest || body.Code != "invalid_argument" {
		t.Fatalf("From(local error) = %d, %#v", code, body)
	}
}

func TestFromHidesInternalDetails(t *testing.T) {
	code, body := From(status.Error(codes.Internal, "postgres password leaked"))
	if code != http.StatusInternalServerError || body.Detail != "internal server error" {
		t.Fatalf("From(internal) = %d, %#v", code, body)
	}
}
