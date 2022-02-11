package entity

import (
	"github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"
	entity "go_scafold/service/transport/endpoints/entity"
	common "go_scafold/service/transport/http/common"
)

func AddHTTPRoutes(r *mux.Router, endpoints interface{}, options []kithttp.ServerOption,
	tracer stdopentracing.Tracer, logger log.Logger) *mux.Route {
	return r.Methods("GET").Path("/entity/{id}").Handler(kithttp.NewServer(
		endpoints.(entity.Endpoints).GetEntity,
		decodeGetEntityRequest,
		common.EncodeResponse,
		append(options, kithttp.ServerBefore(opentracing.HTTPToContext(tracer, "entity", logger)))...))
}
