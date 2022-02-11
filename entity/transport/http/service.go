package http

import (
	"context"
	kittransport "github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go_scafold/entity/transport"
	"net/http"
)

func NewHTTPServer(ctx context.Context, endpoints transport.Endpoints, tracer stdopentracing.Tracer,
	options []kithttp.ServerOption, logger log.Logger) http.Handler {

	r, options := initMuxRouter(logger, options)

	addHTTPRoutes(r, endpoints, options, tracer, logger)

	return r
}

func initMuxRouter(logger log.Logger, options []kithttp.ServerOption) (*mux.Router, []kithttp.ServerOption) {
	r := mux.NewRouter()
	errorEncoder := kithttp.ServerErrorEncoder(encodeErrorResponse)
	errorHandler := kithttp.ServerErrorHandler(kittransport.NewLogErrorHandler(logger))
	options = append(options, errorEncoder, errorHandler)

	r.Use(genericMiddlewareToSetHTTPHeader)
	r.Use(jwtMiddlewareForMicrosoftIdentity)

	r.Methods("GET").Path("/metrics").Handler(promhttp.Handler())
	return r, options
}
