package domain

import (
	"context"
	"time"
)

type FilmService interface {
	Ping(context.Context) error
	CreateFilm(ctx context.Context, title string, description string, releaseDate time.Time, rating float32, actors []int) (int, error)
	UpdateFilm(ctx context.Context, id int, title string, description string, releaseDate time.Time, rating *float32, actors []int) error
	DeleteFilm(ctx context.Context, id int) error
	ReadFilms(ctx context.Context, field string, order string, page int, limit int) ([]OutputFilm, error)
}

type FilmRepository interface {
	Ping(context.Context) error
	CreateFilm(ctx context.Context, title string, description string, releaseDate time.Time, rating float32, actors []int) (int, error)
	UpdateFilm(ctx context.Context, id int, title string, description string, releaseDate time.Time, rating *float32, actors []int) error
	DeleteFilm(ctx context.Context, id int) error
	ReadFilms(ctx context.Context, field string, order string, page int, limit int) ([]OutputFilm, error)
}

type ActorService interface {
	CreateActor(ctx context.Context, name string, gender bool, birthday time.Time) (int, error)
	UpdateActor(ctx context.Context, id int, name string, gender *bool, birthday time.Time) error
	DeleteActor(ctx context.Context, id int) error
}

type ActorRepository interface {
	CreateActor(ctx context.Context, name string, gender bool, birthday time.Time) (int, error)
	UpdateActor(ctx context.Context, id int, name string, gender *bool, birthday time.Time) error
	DeleteActor(ctx context.Context, id int) error
}
