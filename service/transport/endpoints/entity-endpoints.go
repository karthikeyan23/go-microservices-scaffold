package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
	"go_scafold/service/model"
	"time"
)

func getEntityEndpoints(s model.EntityService, logger log.Logger, duration metrics.Histogram, tracer stdopentracing.Tracer,
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

func makeGetEntityEndpoint(s model.EntityService) endpoint.Endpoint {
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
