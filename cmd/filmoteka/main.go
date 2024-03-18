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

	aur := repository.NewAuthorization(repository.NewPostgres(pool))
	ar := repository.NewActor(repository.NewPostgres(pool))
	fr := repository.NewFilm(repository.NewPostgres(pool))
	aus := service.NewAuthorization(aur)
	as := service.NewActor(ar)
	fs := service.NewFilm(fr)
	auh := handlers.NewAuthorization(aus, cfg.JWTKey)
	ah := handlers.NewActor(as)
	fh := handlers.NewFilm(fs)

	mux := http.NewServeMux()

	mux.Handle("GET /ping", middleware.Log(middleware.AuthorizationRequired(http.HandlerFunc(fh.Ping), auh.JWTKey)))
	mux.Handle("POST /actor", middleware.Log(middleware.AdminRequired(http.HandlerFunc(ah.CreateActor), auh.JWTKey)))
	mux.Handle("PUT /actor/{id}", middleware.Log(middleware.AdminRequired(http.HandlerFunc(ah.UpdateActor), auh.JWTKey)))
	mux.Handle("DELETE /actor/{id}", middleware.Log(middleware.AdminRequired(http.HandlerFunc(ah.DeleteActor), auh.JWTKey)))
	mux.Handle("POST /film", middleware.Log(middleware.AdminRequired(http.HandlerFunc(fh.CreateFilm), auh.JWTKey)))
	mux.Handle("PUT /film/{id}", middleware.Log(middleware.AdminRequired(http.HandlerFunc(fh.UpdateFilm), auh.JWTKey)))
	mux.Handle("DELETE /film/{id}", middleware.Log(middleware.AdminRequired(http.HandlerFunc(fh.DeleteFilm), auh.JWTKey)))
	mux.Handle("GET /films", middleware.Log(middleware.AuthorizationRequired(http.HandlerFunc(fh.ReadFilms), auh.JWTKey)))
	mux.Handle("GET /films/search", middleware.Log(middleware.AuthorizationRequired(http.HandlerFunc(fh.FindFilms), auh.JWTKey)))
	mux.Handle("GET /actors", middleware.Log(middleware.AuthorizationRequired(http.HandlerFunc(ah.ReadActors), auh.JWTKey)))
	mux.Handle("POST /register", middleware.Log(http.HandlerFunc(auh.Register)))
	mux.Handle("POST /login", middleware.Log(http.HandlerFunc(auh.LogIn)))

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
