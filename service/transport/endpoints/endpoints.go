package endpoints

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
	domain "go_scafold/service/domain/entity"
	entity "go_scafold/service/transport/endpoints/entity"
)

func MakeEndpoints(s domain.Service, logger log.Logger, duration metrics.Histogram, tracer stdopentracing.Tracer,
) interface{} {

	endpoints := entity.GetEntityEndpoints(s, logger, duration, tracer)
	//append more endpoints here
	return endpoints
}
