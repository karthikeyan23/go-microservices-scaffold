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
		getEntityEndpoint = initEndpoint(makeGetEntityEndpoint(s),
			"get-entity",
			30*time.Second,
			5,
			time.Second,
			logger,
			duration,
			tracer)
	}
	return Endpoints{
		GetEntity: getEntityEndpoint,
	}
}

func initEndpoint(endpoint endpoint.Endpoint, name string, circuitBreakerTimeout time.Duration,
	rateLimit int,
	rateLimitDuration time.Duration,
	logger log.Logger,
	duration metrics.Histogram,
	tracer stdopentracing.Tracer) endpoint.Endpoint {

	endpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(
		rate.Every(rateLimitDuration), rateLimit))(endpoint)
	endpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    name,
		Timeout: circuitBreakerTimeout,
	}))(endpoint)
	endpoint = opentracing.TraceEndpoint(tracer, "get-entity")(endpoint)
	endpoint = instrumentationMiddleware(duration.With("method", "get-entity"))(endpoint)
	return endpoint
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
