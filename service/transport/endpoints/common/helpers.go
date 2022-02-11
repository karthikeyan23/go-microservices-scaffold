package common

import (
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"time"
)

func InitEndpoint(endpoint endpoint.Endpoint, name string, circuitBreakerTimeout time.Duration,
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
	endpoint = opentracing.TraceEndpoint(tracer, name)(endpoint)
	endpoint = InstrumentationMiddleware(duration.With("method", name))(endpoint)
	return endpoint
}
