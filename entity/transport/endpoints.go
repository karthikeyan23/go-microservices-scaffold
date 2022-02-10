package transport

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"go_scafold/entity"
)

type Endpoints struct {
	GetEntity endpoint.Endpoint
}

func MakeEndpoints(s entity.Service) Endpoints {
	return Endpoints{
		GetEntity: makeGetEntityEndpoint(s),
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
