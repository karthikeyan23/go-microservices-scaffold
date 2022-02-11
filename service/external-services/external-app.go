package external_services

import (
	"context"
	"github.com/go-kit/log"
	domain "go_scafold/service/domain/entity"
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
