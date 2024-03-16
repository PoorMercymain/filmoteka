package domain

import (
	"context"
)

type FilmotekaService interface {
	Ping(context.Context) error
}

type FilmotekaRepository interface {
	Ping(context.Context) error
}
