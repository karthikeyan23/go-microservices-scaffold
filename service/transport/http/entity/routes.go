package entity

import (
	"github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"
	endpointcommon "go_scafold/service/transport/endpoints/common"
	"go_scafold/service/transport/http/common"
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
