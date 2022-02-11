package endpoints

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
	domain "go_scafold/service/domain/entity"
)

type Endpoints struct {
	GetEntity endpoint.Endpoint
}

func MakeEndpoints(s domain.Service, logger log.Logger, duration metrics.Histogram, tracer stdopentracing.Tracer,
) Endpoints {

	endpoints := getEntityEndpoints(s, logger, duration, tracer)
	//append more endpoints here
	return endpoints
}
