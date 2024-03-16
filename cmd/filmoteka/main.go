package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/caarlos0/env/v6"

	"github.com/PoorMercymain/filmoteka/internal/filmoteka/config"
	"github.com/PoorMercymain/filmoteka/internal/filmoteka/handlers"
	"github.com/PoorMercymain/filmoteka/internal/filmoteka/middleware"
	"github.com/PoorMercymain/filmoteka/internal/filmoteka/repository"
	"github.com/PoorMercymain/filmoteka/internal/filmoteka/service"
	"github.com/PoorMercymain/filmoteka/pkg/logger"
)

func main() {
	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		logger.Logger().Fatalln("Failed to parse env: %v", err) // default logfile will be used
	}

	logger.SetLogFile(cfg.LogFilePath)

	err := repository.ApplyMigrations(cfg.MigrationsPath, cfg.DSN())
	if err != nil {
		logger.Logger().Fatalln(zap.Error(err))
	}

	logger.Logger().Infoln("Migrations applied successfully")

	pool, err := repository.GetPgxPool(cfg.DSN())
	if err != nil {
		logger.Logger().Fatalln(zap.Error(err))
	}

	logger.Logger().Infoln("Postgres connection pool created")

	r := repository.New(repository.NewPostgres(pool))
	s := service.New(r)
	h := handlers.New(s)

	mux := http.NewServeMux()

	mux.Handle("GET /ping", middleware.Log(http.HandlerFunc(h.Ping)))

	server := &http.Server{
		Addr:     cfg.ServiceHost + ":" + strconv.Itoa(cfg.ServicePort),
		ErrorLog: log.New(logger.Logger(), "", 0),
		Handler:  mux,
	}

	go func() {
		logger.Logger().Infoln("Server started, listening on port", cfg.ServicePort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger().Fatalln("ListenAndServe failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	logger.Logger().Infoln("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Logger().Fatalln("Server was forced to shutdown:", zap.Error(err))
	}

	logger.Logger().Infoln("Server was shut down")
}
