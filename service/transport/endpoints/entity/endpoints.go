package entity

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/log"
	domain "github.com/karthkeyan23/go_microservices_scaffold/service/domain/entity"
	common "github.com/karthkeyan23/go_microservices_scaffold/service/transport/endpoints/common"
	stdopentracing "github.com/opentracing/opentracing-go"
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

	var getAppDataEndPoint endpoint.Endpoint
	{
		getAppDataEndPoint = common.InitEndpoint(makeGetAppDataEndpoint(s),
			"get-app-data",
			30*time.Second,
			5,
			time.Second,
			logger,
			duration,
			tracer)
	}
	endpoints.GetAppData = getAppDataEndPoint
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

func makeGetAppDataEndpoint(s domain.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetAppDataRequest)
		data, err := s.GetDataFromApp(ctx, req)
		var res GetAppDataResponse
		if err != nil {
			return nil, err
		}
		res.Data = data.(bool)
		return res, nil
	}
}
