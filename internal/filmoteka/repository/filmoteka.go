package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
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

func (r *filmoteka) CreateActor(ctx context.Context, name string, gender bool, birthday time.Time) (int, error) {
	var id int
	err := r.db.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		return tx.QueryRow(ctx, "INSERT INTO actors(name, gender, birthday) VALUES($1, $2, $3) RETURNING id", name, gender, birthday).Scan(&id)
	})

	if err != nil {
		return 0, fmt.Errorf("repository.CreateActor(): %w", err)
	}

	return id, nil
}
