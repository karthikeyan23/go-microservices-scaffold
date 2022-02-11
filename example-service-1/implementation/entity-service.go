package implementation

import (
	"context"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"go_scafold/example-service-1/model"
)

// The Business service-implementation of the Service service-interface.
type service struct {
	repository model.Repository
	logger     log.Logger
}

// NewService Creates the service and returns a pointer with Service methods implemented.
func NewService(repository model.Repository, logger log.Logger) model.EntityService {
	return &service{
		repository: repository,
		logger:     logger,
	}
}

func (s service) GetEntity(ctx context.Context, id string) (*model.Entity, error) {
	logger := log.With(s.logger, "method", "get-entity")
	aEntity, err := s.repository.Get(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		return nil, err
	}
	if aEntity == nil {
		return nil, model.ErrEntityNotFound
	}
	return aEntity, nil
}

func (s service) CreateEntity(ctx context.Context, entity *model.Entity) error {
	//TODO implement me
	panic("implement me")
}

func (s service) UpdateEntity(ctx context.Context, user *model.Entity) error {
	//TODO implement me
	panic("implement me")
}

func (s service) DeleteEntity(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
