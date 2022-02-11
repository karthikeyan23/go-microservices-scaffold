package entrypoints

import (
	"database/sql"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
	common "go_scafold/service/transport/endpoints/common"
)

func InitServicesAndEndPoints(db *sql.DB, logger log.Logger, duration metrics.Histogram,
	tracer stdopentracing.Tracer) (common.Endpoints, error) {
	var endpoints common.Endpoints
	addEntityServicesAndGetEndpoints(&endpoints, db, logger, duration, tracer)
	//Add more services here
	return endpoints, nil
}
