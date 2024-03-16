package repository

import (
	"context"
	"fmt"
)

type filmoteka struct {
	db *postgres
}

func New(pg *postgres) *filmoteka {
	return &filmoteka{db: pg}
}

func (r *filmoteka) Ping(ctx context.Context) error {
	err := r.db.Ping(ctx)
	if err != nil {
		return fmt.Errorf("repository.Ping(): %w", err)
	}

	return nil
}
