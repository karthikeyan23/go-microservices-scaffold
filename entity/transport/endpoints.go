package transport

import (
	"context"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
	"go_scafold/entity"
	"golang.org/x/time/rate"
	"time"
)

type Endpoints struct {
	GetEntity endpoint.Endpoint
}

func MakeEndpoints(s entity.Service, logger log.Logger, duration metrics.Histogram, tracer stdopentracing.Tracer,
) Endpoints {

	var getEntityEndpoint endpoint.Endpoint
	{
		getEntityEndpoint = makeGetEntityEndpoint(s)
		getEntityEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(
			rate.Every(time.Second), 10))(getEntityEndpoint)
		getEntityEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "get-entity",
			Timeout: 30 * time.Second,
		}))(getEntityEndpoint)
		getEntityEndpoint = opentracing.TraceEndpoint(tracer, "get-entity")(getEntityEndpoint)
		getEntityEndpoint = instrumentingMiddleware(duration.With("method", "get-entity"))(getEntityEndpoint)
	}
	return Endpoints{
		GetEntity: getEntityEndpoint,
	}
}

func makeGetEntityEndpoint(s entity.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetEntityByIDRequest)
		aEntity, err := s.GetEntity(ctx, req.ID)
		if err != nil {
			return nil, err
		}
		return GetEntityByIDResponse{
				ID:        aEntity.ID,
				Name:      aEntity.Name,
				CreatedAt: aEntity.CreatedAt},
			nil
	}
}
