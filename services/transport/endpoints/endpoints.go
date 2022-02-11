package endpoints

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
	"go_scafold/services/model"
)

type Endpoints struct {
	GetEntity endpoint.Endpoint
}

func MakeEndpoints(s model.EntityService, logger log.Logger, duration metrics.Histogram, tracer stdopentracing.Tracer,
) Endpoints {

	endpoints := getEntityEndpoints(s, logger, duration, tracer)
	//append more endpoints here
	return endpoints
}
