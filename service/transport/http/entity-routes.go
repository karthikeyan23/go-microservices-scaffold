package http

import (
	"github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"
	transport "go_scafold/service/transport/endpoints"
)

func addHTTPRoutes(r *mux.Router, endpoints transport.Endpoints, options []kithttp.ServerOption,
	tracer stdopentracing.Tracer, logger log.Logger) *mux.Route {
	return r.Methods("GET").Path("/entity/{id}").Handler(kithttp.NewServer(
		endpoints.GetEntity,
		decodeGetEntityRequest,
		encodeResponse,
		append(options, kithttp.ServerBefore(opentracing.HTTPToContext(tracer, "entity", logger)))...))
}
