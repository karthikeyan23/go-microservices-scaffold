package http

import (
	"context"
	"encoding/json"
	"errors"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"go_scafold/entity"
	"go_scafold/entity/transport"
	"net/http"
)

var (
	ErrEmptyAuthHeader = errors.New("authorization header is empty")
	ErrMalformedToken  = errors.New("malformed Token")
	ErrBadRouting      = errors.New("bad routing")
)

func NewHTTPServer(endpoints transport.Endpoints, options []kithttp.ServerOption) http.Handler {
	r := mux.NewRouter()
	errorEncoder := kithttp.ServerErrorEncoder(encodeErrorResponse)
	options = append(options, errorEncoder)

	r.Use(genericMiddlewareToSetHTTPHeader)
	r.Use(jwtMiddlewareForMicrosoftIdentity)

	r.Methods("GET").Path("/entity/{id}").Handler(kithttp.NewServer(
		endpoints.GetEntity,
		decodeGetEntityRequest,
		encodeResponse,
		options...))

	return r
}

func decodeGetEntityRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	if err != nil {
		return nil, err
	}
	return transport.GetEntityByIDRequest{ID: id}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeErrorResponse(ctx, e.error(), w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

func encodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case entity.ErrEntityNotFound:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
