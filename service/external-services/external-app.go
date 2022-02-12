package external_services

import (
	"context"
	"github.com/go-kit/log"
	domain "github.com/karthkeyan23/go_microservices_scaffold/service/domain/entity"
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
