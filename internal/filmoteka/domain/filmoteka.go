package domain

import (
	"context"
	"time"
)

type FilmService interface {
	Ping(context.Context) error
	CreateFilm(ctx context.Context, title string, description string, releaseDate time.Time, rating int, actors []int) (int, error)
}

type FilmRepository interface {
	Ping(context.Context) error
	CreateFilm(ctx context.Context, title string, description string, releaseDate time.Time, rating int, actors []int) (int, error)
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
