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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go_scafold/entity"
	repo "go_scafold/entity/db"
	service "go_scafold/entity/implementation"
	transport "go_scafold/entity/transport"
	httptransport "go_scafold/entity/transport/http"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//Load command line flags
	var httpAddr = flag.String("http", ":8080", "HTTP listen address")
	flag.Parse()
	//Load environment variables for the Database connection string
	dbSource := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("SSL_MODE"))

	//Create a logger
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
		return
	}
	//Print the log on service exit
	defer func(info log.Logger, keyvals ...interface{}) {
		err := info.Log(keyvals)
		if err != nil {

		}
	}(level.Info(logger), "msg", "service ended")
	//Add OpenTracing tacker
	tracer := stdopentracing.GlobalTracer()
	//Create sparse metrics
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
	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())

	ctx := context.Background()
	//Initialize the DataBase
	var db *sql.DB
	{
		db, err = sql.Open("postgres", dbSource)
		if err != nil {
			err := level.Error(logger).Log("exit", err)
			if err != nil {
				return
			}
			os.Exit(-1)
		}
	}
	//Close the database connection on service exit
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			err := level.Error(logger).Log("exit", err)
			if err != nil {
				return
			}
			os.Exit(-1)
		}
	}(db)
	//Initialize the repository
	var svc entity.Service
	{
		repository, err := repo.New(db, logger)
		if err != nil {
			err := level.Error(logger).Log("exit", err)
			if err != nil {
				return
			}
			os.Exit(-1)
		}
		svc = service.NewService(repository, logger)
	}
	//Initialize the Endpoints
	endpoints := transport.MakeEndpoints(svc, logger, duration, tracer)
	//initialize the HTTP transport
	var h http.Handler
	{
		var serverOptions []kithttp.ServerOption
		h = httptransport.NewHTTPServer(ctx, endpoints, tracer, serverOptions, logger)
	}
	//Channel to listen for service exit
	errChannel := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChannel <- fmt.Errorf("%s", <-c)
	}()
	//Start the HTTP server
	go func() {
		err := level.Info(logger).Log("transport", "HTTP", "addr", *httpAddr)
		if err != nil {
			return
		}
		server := &http.Server{
			Addr:    *httpAddr,
			Handler: h,
		}
		errChannel <- server.ListenAndServe()
	}()
	//Print the error on service exit
	err = level.Error(logger).Log("exit", <-errChannel)
	if err != nil {
		return
	}
}
