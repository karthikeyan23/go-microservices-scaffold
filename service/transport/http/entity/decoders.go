package entity

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	entity "github.com/karthkeyan23/go_microservices_scaffold/service/transport/endpoints/entity"
	common "github.com/karthkeyan23/go_microservices_scaffold/service/transport/http/common"
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

func decodeGeAppDataRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req entity.GetAppDataRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}
