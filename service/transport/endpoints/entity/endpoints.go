package entity

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
	domain "go_scafold/service/domain/entity"
	common "go_scafold/service/transport/endpoints/common"
	"time"
)

func MakeEndpoints(endpoints *common.Endpoints, s domain.Service, logger log.Logger, duration metrics.Histogram,
	tracer stdopentracing.Tracer) {
	var getEntityEndpoint endpoint.Endpoint
	{
		getEntityEndpoint = common.InitEndpoint(makeGetEntityEndpoint(s),
			"get-entity",
			30*time.Second,
			5,
			time.Second,
			logger,
			duration,
			tracer)
	}
	endpoints.GetEntity = getEntityEndpoint
}

func makeGetEntityEndpoint(s domain.Service) endpoint.Endpoint {
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
