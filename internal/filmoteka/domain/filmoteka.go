package domain

import (
	"context"
	"time"
)

type FilmService interface {
	CreateFilm(ctx context.Context, title string, description string, releaseDate time.Time, rating float32, actors []int) (int, error)
	UpdateFilm(ctx context.Context, id int, title string, description string, releaseDate time.Time, rating *float32, actors []int) error
	DeleteFilm(ctx context.Context, id int) error
	ReadFilms(ctx context.Context, field string, order string, page int, limit int) ([]OutputFilm, error)
	FindFilms(ctx context.Context, filmTitleFragment string, actorNameFragment string, page int, limit int) ([]OutputFilm, error)
}

//go:generate mockgen -destination=mocks/film_repo_mock.gen.go -package=mocks . FilmRepository
type FilmRepository interface {
	CreateFilm(ctx context.Context, title string, description string, releaseDate time.Time, rating float32, actors []int) (int, error)
	UpdateFilm(ctx context.Context, id int, title string, description string, releaseDate time.Time, rating *float32, actors []int) error
	DeleteFilm(ctx context.Context, id int) error
	ReadFilms(ctx context.Context, field string, order string, page int, limit int) ([]OutputFilm, error)
	FindFilms(ctx context.Context, filmTitleFragment string, actorNameFragment string, page int, limit int) ([]OutputFilm, error)
}

type ActorService interface {
	CreateActor(ctx context.Context, name string, gender bool, birthday time.Time) (int, error)
	UpdateActor(ctx context.Context, id int, name string, gender *bool, birthday time.Time) error
	DeleteActor(ctx context.Context, id int) error
	ReadActors(ctx context.Context, page int, limit int) ([]OutputActor, error)
}

//go:generate mockgen -destination=mocks/actor_repo_mock.gen.go -package=mocks . ActorRepository
type ActorRepository interface {
	CreateActor(ctx context.Context, name string, gender bool, birthday time.Time) (int, error)
	UpdateActor(ctx context.Context, id int, name string, gender *bool, birthday time.Time) error
	DeleteActor(ctx context.Context, id int) error
	ReadActors(ctx context.Context, page int, limit int) ([]OutputActor, error)
}

type AuthorizationService interface {
	Register(ctx context.Context, login string, password string) error
	CheckAuth(ctx context.Context, login string, password string) error
	IsAdmin(ctx context.Context, login string) (bool, error)
}

//go:generate mockgen -destination=mocks/authorization_repo_mock.gen.go -package=mocks . AuthorizationRepository
type AuthorizationRepository interface {
	Register(ctx context.Context, login string, passwordHash string) error
	GetPasswordHash(ctx context.Context, login string) (string, error)
	IsAdmin(ctx context.Context, login string) (bool, error)
}
