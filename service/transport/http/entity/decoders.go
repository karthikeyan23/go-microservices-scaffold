package entity

import (
	"context"
	"github.com/gorilla/mux"
	entity "go_scafold/service/transport/endpoints/entity"
	common "go_scafold/service/transport/http/common"
	"net/http"
)

func decodeGetEntityRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, common.ErrBadRouting
	}
	if err != nil {
		return nil, err
	}
	return entity.GetEntityByIDRequest{ID: id}, nil
}
