package http

import (
	"context"
	kittransport "github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	common "go_scafold/service/transport/http/common"
	"go_scafold/service/transport/http/entity"
	"net/http"
)

func NewHTTPServer(ctx context.Context, endpoints interface{}, tracer stdopentracing.Tracer,
	options []kithttp.ServerOption, logger log.Logger) http.Handler {
	r, options := initMuxRouter(logger, options)
	entity.AddHTTPRoutes(r, endpoints, options, tracer, logger)
	//append more routes here
	return r
}

func initMuxRouter(logger log.Logger, options []kithttp.ServerOption) (*mux.Router, []kithttp.ServerOption) {
	r := mux.NewRouter()
	errorEncoder := kithttp.ServerErrorEncoder(common.EncodeErrorResponse)
	errorHandler := kithttp.ServerErrorHandler(kittransport.NewLogErrorHandler(logger))
	options = append(options, errorEncoder, errorHandler)

	r.Use(common.GenericMiddlewareToSetHTTPHeader)
	r.Use(common.JwtMiddlewareForMicrosoftIdentity)

	r.Methods("GET").Path("/metrics").Handler(promhttp.Handler())
	return r, options
}
