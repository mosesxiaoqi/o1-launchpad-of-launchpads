package problem

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Body struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
	Code   string `json:"code"`
}

func From(err error) (int, Body) {
	grpcStatus, isGRPC := status.FromError(err)
	grpcCode := grpcStatus.Code()
	httpStatus, code := http.StatusInternalServerError, "internal_error"
	if !isGRPC {
		httpStatus, code = http.StatusBadRequest, "invalid_argument"
	}
	switch grpcCode {
	case codes.InvalidArgument:
		httpStatus, code = http.StatusBadRequest, "invalid_argument"
	case codes.Unauthenticated:
		httpStatus, code = http.StatusUnauthorized, "unauthenticated"
	case codes.PermissionDenied:
		httpStatus, code = http.StatusForbidden, "permission_denied"
	case codes.NotFound:
		httpStatus, code = http.StatusNotFound, "not_found"
	case codes.AlreadyExists, codes.FailedPrecondition, codes.Aborted:
		httpStatus, code = http.StatusConflict, "conflict"
	case codes.DeadlineExceeded:
		httpStatus, code = http.StatusGatewayTimeout, "deadline_exceeded"
	case codes.Unavailable:
		httpStatus, code = http.StatusServiceUnavailable, "unavailable"
	}
	detail := grpcStatus.Message()
	if httpStatus == http.StatusInternalServerError {
		detail = "internal server error"
	}
	return httpStatus, Body{
		Type: "about:blank", Title: http.StatusText(httpStatus), Status: httpStatus, Detail: detail, Code: code,
	}
}
