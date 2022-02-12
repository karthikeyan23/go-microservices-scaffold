package entrypoints

import (
	"database/sql"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	repo "github.com/karthkeyan23/go_microservices_scaffold/service/db"
	domain "github.com/karthkeyan23/go_microservices_scaffold/service/domain/entity"
	app "github.com/karthkeyan23/go_microservices_scaffold/service/external-services"
	common "github.com/karthkeyan23/go_microservices_scaffold/service/transport/endpoints/common"
	entity "github.com/karthkeyan23/go_microservices_scaffold/service/transport/endpoints/entity"
	stdopentracing "github.com/opentracing/opentracing-go"
	"os"
)

func addEntityServicesAndGetEndpoints(endpoints *common.Endpoints, db *sql.DB, logger log.Logger,
	duration metrics.Histogram, tracer stdopentracing.Tracer) {
	//Initialize the entity repository
	svc := initRepoAndService(db, logger)
	//Initialize the entity Endpoints
	entity.MakeEndpoints(endpoints, svc, logger, duration, tracer)
}

func initRepoAndService(db *sql.DB, logger log.Logger) domain.Service {
	var svc domain.Service
	{
		repository, err := repo.New(db, logger)
		if err != nil {
			_ = level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
		externalApp, err := app.NewExternalApp(logger)
		if err != nil {
			_ = level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
		svc = domain.NewService(repository, externalApp, logger)
	}
	return svc
}
