package http

import (
	"context"
	"github.com/gorilla/mux"
	transport "go_scafold/service/transport/endpoints"
	"net/http"
)

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
