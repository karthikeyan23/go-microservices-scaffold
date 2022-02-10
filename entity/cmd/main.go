package main

import (
	"database/sql"
	"flag"
	"fmt"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	_ "github.com/lib/pq"
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

var dbSource = "postgresql://beezlabs:plugmein@localhost:5432/kn?sslmode=disable"

func main() {
	var httpAddr = flag.String("http", ":8080", "HTTP listen address")

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		logger = log.With(logger,
			"service", "entity",
			"time:", log.DefaultTimestampUTC,
			"caller:", log.DefaultCaller,
		)
	}
	err := level.Info(logger).Log("msg", "service started")
	if err != nil {
		return
	}
	defer func(info log.Logger, keyvals ...interface{}) {
		err := info.Log(keyvals)
		if err != nil {

		}
	}(level.Info(logger), "msg", "service ended")

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

	flag.Parse()

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

	endpoints := transport.MakeEndpoints(svc)

	var h http.Handler
	{
		var serverOptions []kithttp.ServerOption
		h = httptransport.NewHTTPServer(endpoints, serverOptions)
	}

	errChannel := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChannel <- fmt.Errorf("%s", <-c)
	}()

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
	err = level.Error(logger).Log("exit", <-errChannel)
	if err != nil {
		return
	}
}
