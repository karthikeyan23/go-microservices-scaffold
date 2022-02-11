package entrypoints

import (
	"database/sql"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	stdopentracing "github.com/opentracing/opentracing-go"
	repo "go_scafold/service/db"
	domain "go_scafold/service/domain/entity"
	common "go_scafold/service/transport/endpoints/common"
	entity "go_scafold/service/transport/endpoints/entity"
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
		svc = domain.NewService(repository, logger)
	}
	return svc
}
