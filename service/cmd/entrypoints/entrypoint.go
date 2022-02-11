package entrypoints

import (
	"database/sql"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
	transport "go_scafold/service/transport/endpoints"
)

func InitServicesAndEndPoints(db *sql.DB, logger log.Logger, duration metrics.Histogram,
	tracer stdopentracing.Tracer) transport.Endpoints {
	endpoints := addEntityServicesAndGetEndpoints(db, logger, duration, tracer)
	//Add more services here
	return endpoints
}
