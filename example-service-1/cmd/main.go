package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	_ "github.com/lib/pq"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	repo "go_scafold/example-service-1/db"
	service "go_scafold/example-service-1/implementation"
	"go_scafold/example-service-1/model"
	transport "go_scafold/example-service-1/transport/endpoints"
	httptransport "go_scafold/example-service-1/transport/http"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//Load command line flags
	httpAddr := setApplicationPort()
	//Load environment variables for the Database connection string
	dbSource := setDBConnectionString()

	//Create a logger
	logger, err := initLogger()
	if err != nil {
		return
	}
	//Print the log on service exit
	defer onServiceClose(logger)
	//Add OpenTracing tacker
	tracer := stdopentracing.GlobalTracer()
	//Create sparse metrics
	duration := initDurationMetrics()
	//Create a context
	ctx := context.Background()

	db := initDB(dbSource, logger)
	//Close the database connection on service exit
	defer closeDB(db, logger)
	//Initialise all services in the project
	endpoints := initServicesAndEndPoints(db, logger, duration, tracer)
	//initialize the HTTP transport
	httpTransportHandler := addHTTPTransport(ctx, endpoints, tracer, logger)
	//Channel to listen for service exit
	errChannel := make(chan error)
	go waitForInterrupt(errChannel)
	//Start the HTTP server
	go startHttpServer(logger, httpTransportHandler, httpAddr, errChannel)
	//Print the error on service exit
	_ = level.Error(logger).Log("exit", <-errChannel)
}

func initServicesAndEndPoints(db *sql.DB, logger log.Logger, duration metrics.Histogram,
	tracer stdopentracing.Tracer) transport.Endpoints {
	endpoints := addEntityServicesAndGetEndpoints(db, logger, duration, tracer)
	//Add more services here
	return endpoints
}

func addEntityServicesAndGetEndpoints(db *sql.DB, logger log.Logger, duration metrics.Histogram,
	tracer stdopentracing.Tracer) transport.Endpoints {
	//Initialize the entity repository
	svc := initRepoAndService(db, logger)
	//Initialize the entity Endpoints
	endpoints := transport.MakeEndpoints(svc, logger, duration, tracer)
	return endpoints
}

func addHTTPTransport(ctx context.Context, endpoints transport.Endpoints, tracer stdopentracing.Tracer, logger log.Logger) http.Handler {
	var h http.Handler
	{
		var serverOptions []kithttp.ServerOption
		h = httptransport.NewHTTPServer(ctx, endpoints, tracer, serverOptions, logger)
	}
	return h
}

func onServiceClose(logger log.Logger) {
	_ = level.Info(logger).Log("msg", "service terminating")
}
func closeDB(db *sql.DB, logger log.Logger) {
	err := db.Close()
	if err != nil {
		_ = level.Error(logger).Log("exit", err)
		os.Exit(-1)
	}
	_ = level.Info(logger).Log("msg", "database connection closed")
}

func startHttpServer(logger log.Logger, h http.Handler, httpAddr *string, errChannel chan error) {
	err := level.Info(logger).Log("transport", "HTTP", "addr", *httpAddr)
	if err != nil {
		return
	}
	server := &http.Server{
		Addr:    *httpAddr,
		Handler: h,
	}
	errChannel <- server.ListenAndServe()
}

func waitForInterrupt(errChannel chan error) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	errChannel <- fmt.Errorf("%s", <-c)
}

func initRepoAndService(db *sql.DB, logger log.Logger) model.EntityService {
	var svc model.EntityService
	{
		repository, err := repo.New(db, logger)
		if err != nil {
			_ = level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
		svc = service.NewService(repository, logger)
	}
	return svc
}

func initDB(dbSource string, logger log.Logger) *sql.DB {
	var err error
	var db *sql.DB
	{
		db, err = sql.Open("postgres", dbSource)
		if err != nil {
			_ = level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
	}
	return db
}

func initDurationMetrics() metrics.Histogram {
	var duration metrics.Histogram
	{
		// Endpoint-level metrics.
		duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "beezlabs",
			Subsystem: "entity_service",
			Name:      "request_duration_seconds",
			Help:      "Request duration in seconds.",
		}, []string{"method", "success"})
	}
	return duration
}

func initLogger() (log.Logger, error) {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		logger = log.With(logger,
			"service", "entity",
			"time:", log.DefaultTimestampUTC,
			"caller:", log.DefaultCaller,
		)
	}
	//Add the first log for the service
	err := level.Info(logger).Log("msg", "service started")
	if err != nil {
		return nil, nil
	}
	return logger, err
}

func setDBConnectionString() string {
	dbSource := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("SSL_MODE"))
	return dbSource
}

func setApplicationPort() *string {
	var httpAddr = flag.String("http", ":8080", "HTTP listen address")
	flag.Parse()
	return httpAddr
}
