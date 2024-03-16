package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/PoorMercymain/filmoteka/internal/filmoteka/handlers"
	"github.com/PoorMercymain/filmoteka/internal/filmoteka/middleware"
	"github.com/PoorMercymain/filmoteka/internal/filmoteka/repository"
	"github.com/PoorMercymain/filmoteka/internal/filmoteka/service"
	"github.com/PoorMercymain/filmoteka/pkg/logger"
)

func main() {
	migrationsPath := "migrations"
	dsn := "postgres://filmoteka:filmoteka@localhost:5432/filmoteka?sslmode=disable"

	err := repository.ApplyMigrations(migrationsPath, dsn)
	if err != nil {
		logger.Logger().Fatalln(zap.Error(err))
		return
	}

	logger.Logger().Infoln("Migrations applied successfully")

	pool, err := repository.GetPgxPool(dsn)
	if err != nil {
		logger.Logger().Fatalln(zap.Error(err))
		return
	}

	logger.Logger().Infoln("Postgres connection pool created")

	r := repository.New(repository.NewPostgres(pool))
	s := service.New(r)
	h := handlers.New(s)

	mux := http.NewServeMux()

	mux.Handle("GET /ping", middleware.Log(http.HandlerFunc(h.Ping)))

	server := &http.Server{
		Addr:     "localhost:8080",
		ErrorLog: log.New(logger.Logger(), "", 0),
		Handler:  mux,
	}

	go func() {
		logger.Logger().Infoln("Server started, listening on port", 8080)
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
