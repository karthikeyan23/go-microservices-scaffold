package common

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/ratelimit"
	"github.com/sony/gobreaker"
	domain "go_scafold/service/domain/entity"
	"net/http"
)

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorInterface); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		EncodeErrorResponse(ctx, e.error(), w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}

type errorInterface interface {
	error() error
}

func EncodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case gobreaker.ErrTooManyRequests:
		return http.StatusTooManyRequests
	case gobreaker.ErrOpenState:
		return http.StatusServiceUnavailable
	case domain.ErrEntityNotFound:
		return http.StatusBadRequest
	case ratelimit.ErrLimited:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
