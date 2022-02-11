package entity

import (
	"context"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// The Business service-implementation of the Service service-interface.
type service struct {
	externalApp ExternalApp
	repository  Repository
	logger      log.Logger
}

// NewService Creates the service and returns a pointer with Service methods implemented.
func NewService(repository Repository, externalApp ExternalApp, logger log.Logger) Service {
	return &service{
		externalApp: externalApp,
		repository:  repository,
		logger:      logger,
	}
}

func (s service) GetEntity(ctx context.Context, id string) (*Entity, error) {
	logger := log.With(s.logger, "method", "get-entity")
	aEntity, err := s.repository.Get(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		return nil, err
	}
	if aEntity == nil {
		return nil, ErrEntityNotFound
	}
	return aEntity, nil
}

func (s service) CreateEntity(ctx context.Context, entity *Entity) error {
	//TODO implement me
	panic("implement me")
}

func (s service) UpdateEntity(ctx context.Context, user *Entity) error {
	//TODO implement me
	panic("implement me")
}

func (s service) DeleteEntity(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (s service) GetDataFromApp(ctx context.Context, input interface{}) (interface{}, error) {
	return s.externalApp.GetData(ctx, input)
}
