package external_services

import (
	"context"
	domain "github.com/beezlabs-org/go_microservices_scaffold/service/domain/entity"
	"github.com/go-kit/log"
)

type externalApp struct {
	logger log.Logger
}

func NewExternalApp(logger log.Logger) (domain.ExternalApp, error) {
	return externalApp{logger: logger}, nil
}

func (e externalApp) GetData(ctx context.Context, input interface{}) (interface{}, error) {
	return true, nil
}
