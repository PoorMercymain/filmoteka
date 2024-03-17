package domain

import (
	"context"
	"time"
)

type FilmotekaService interface {
	Ping(context.Context) error
	CreateActor(ctx context.Context, name string, gender bool, birthday time.Time) (int, error)
}

type FilmotekaRepository interface {
	Ping(context.Context) error
	CreateActor(ctx context.Context, name string, gender bool, birthday time.Time) (int, error)
}
