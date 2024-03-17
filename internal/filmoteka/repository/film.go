package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	appErrors "github.com/PoorMercymain/filmoteka/errors"
	"github.com/PoorMercymain/filmoteka/internal/filmoteka/domain"
)

var (
	_ domain.FilmRepository = (*film)(nil)
)

type film struct {
	db *postgres
}

func NewFilm(pg *postgres) *film {
	return &film{db: pg}
}

func (r *film) Ping(ctx context.Context) error {
	err := r.db.Ping(ctx)
	if err != nil {
		return fmt.Errorf("repository.Ping(): %w", err)
	}

	return nil
}

func (r *film) CreateFilm(ctx context.Context, title string, description string, releaseDate time.Time, rating int, actors []int) (int, error) {
	var id int
	err := r.db.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		err := tx.QueryRow(ctx, "INSERT INTO films(title, description, release_date, rating) VALUES($1, $2, $3, $4) RETURNING id", title, description, releaseDate, rating).Scan(&id)
		if err != nil {
			return err
		}

		for _, actorID := range actors {
			_, err = tx.Exec(ctx, "INSERT INTO film_actor(actor_id, film_id) VALUES($1, $2)", actorID, id)
			if err != nil {
				break
			}
		}

		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == "P0001" {
					return appErrors.ErrActorNotBornBeforeFilmRelease
				}

				if pgErr.Code == pgerrcode.ForeignKeyViolation {
					return appErrors.ErrActorDoesNotExist
				}
			}

			return err
		}

		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("repository.CreateFilm(): %w", err)
	}

	return id, nil
}
