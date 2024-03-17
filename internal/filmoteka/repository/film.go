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

func (r *film) UpdateFilm(ctx context.Context, id int, title string, description string, releaseDate time.Time, rating *int, actors []int) error {
	var (
		titleInDB       string
		descriptionInDB string
		releaseDateInDB time.Time
		ratingInDB      int
		actorsInDB      []int
	)

	err := r.db.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		err := tx.QueryRow(ctx, "SELECT title, description, release_date, rating FROM films WHERE id = $1", id).Scan(&titleInDB, &descriptionInDB, &releaseDateInDB, &ratingInDB)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return appErrors.ErrNotFoundInDB
			}

			return err
		}

		if actors == nil {
			rows, err := tx.Query(ctx, "SELECT actor_id FROM film_actor WHERE film_id = $1", id)
			if err != nil {
				return err
			}

			var actorID int
			for rows.Next() {
				err = rows.Scan(&actorID)
				if err != nil {
					return err
				}

				actorsInDB = append(actorsInDB, actorID)
			}

			actors = actorsInDB
		}

		if title == "" {
			title = titleInDB
		}

		if description == "" {
			description = descriptionInDB
		}

		var defaultTime time.Time
		if releaseDate == defaultTime {
			releaseDate = releaseDateInDB
		}

		if rating == nil {
			rating = &ratingInDB
		}

		tag, err := tx.Exec(ctx, "UPDATE films SET title = $1, description = $2, release_date = $3, rating = $4 WHERE id = $5", title, description, releaseDate, *rating, id)
		if err != nil {
			return err
		}

		if tag.RowsAffected() == 0 {
			return appErrors.ErrNotFoundInDB
		}

		_, err = tx.Exec(ctx, "DELETE FROM film_actor WHERE film_id = $1", id)
		if err != nil {
			return err
		}

		for _, actorID := range actors {
			_, err = tx.Exec(ctx, "INSERT INTO film_actor(actor_id, film_id) VALUES($1, $2)", actorID, id)
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
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("repository.UpdateFilm(): %w", err)
	}

	return nil
}

func (r *film) DeleteFilm(ctx context.Context, id int) error {
	err := r.db.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, "DELETE FROM films WHERE id = $1", id)
		if err != nil {
			return err
		}

		if tag.RowsAffected() == 0 {
			return appErrors.ErrNotFoundInDB
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("repository.DeleteFilm(): %w", err)
	}

	return nil
}
