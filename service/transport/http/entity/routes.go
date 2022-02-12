package entity

import (
	endpointcommon "github.com/beezlabs-org/go_microservices_scaffold/service/transport/endpoints/common"
	"github.com/beezlabs-org/go_microservices_scaffold/service/transport/http/common"
	"github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"
)

func AddHTTPRoutes(r *mux.Router, endpoints endpointcommon.Endpoints, options []kithttp.ServerOption,
	tracer stdopentracing.Tracer, logger log.Logger) {
	r.Methods("GET").Path("/entity/app").Handler(kithttp.NewServer(
		endpoints.GetAppData,
		decodeGeAppDataRequest,
		common.EncodeResponse,
		append(options, kithttp.ServerBefore(opentracing.HTTPToContext(tracer, "get app data", logger)))...))

	r.Methods("GET").Path("/entity/{id}").Handler(kithttp.NewServer(
		endpoints.GetEntity,
		decodeGetEntityRequest,
		common.EncodeResponse,
		append(options, kithttp.ServerBefore(opentracing.HTTPToContext(tracer, "get entity", logger)))...))
}
