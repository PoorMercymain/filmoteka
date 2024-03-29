package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

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

func (r *film) CreateFilm(ctx context.Context, title string, description string, releaseDate time.Time, rating float32, actors []int) (int, error) {
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

func (r *film) UpdateFilm(ctx context.Context, id int, title string, description string, releaseDate time.Time, rating *float32, actors []int) error {
	var (
		titleInDB       string
		descriptionInDB string
		releaseDateInDB time.Time
		ratingInDB      float32
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

func (r *film) ReadFilms(ctx context.Context, field string, order string, page int, limit int) ([]domain.OutputFilm, error) {
	var films []domain.OutputFilm
	err := r.db.WithConnection(ctx, func(ctx context.Context, c *pgxpool.Conn) error {
		sqlStr := fmt.Sprintf("SELECT id, title, description, release_date, rating FROM films ORDER BY %s %s", field, order)
		rows, err := c.Query(ctx, sqlStr+" LIMIT $1 OFFSET $2", limit, (page-1)*limit)
		if err != nil {
			return err
		}

		var curFilm domain.OutputFilm
		var curReleaseDate time.Time
		for rows.Next() {
			err = rows.Scan(&curFilm.ID, &curFilm.Title, &curFilm.Description, &curReleaseDate, &curFilm.Rating)
			if err != nil {
				return err
			}

			curFilm.ReleaseDate = curReleaseDate.Format(time.DateOnly)

			actorRows, err := r.db.Query(ctx, "SELECT actor_id FROM film_actor WHERE film_id = $1 ORDER BY actor_id ASC", curFilm.ID)
			if err != nil {
				return err
			}

			var (
				curActor    domain.FilmOutputActor
				curGender   bool
				curBirthday time.Time
			)

			curFilm.Actors = make([]domain.FilmOutputActor, 0)
			for actorRows.Next() {
				err = actorRows.Scan(&curActor.ID)
				if err != nil {
					return err
				}

				err = r.db.QueryRow(ctx, "SELECT name, gender, birthday FROM actors WHERE id = $1", curActor.ID).Scan(&curActor.Name, &curGender, &curBirthday)
				if err != nil {
					return err
				}

				if curGender {
					curActor.Gender = "female"
				} else {
					curActor.Gender = "male"
				}

				curActor.Birthday = curBirthday.Format(time.DateOnly)

				curFilm.Actors = append(curFilm.Actors, curActor)
			}

			films = append(films, curFilm)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("repository.ReadFilms(): %w", err)
	}

	return films, nil
}

func (r *film) FindFilms(ctx context.Context, filmTitleFragment string, actorNameFragment string, page int, limit int) ([]domain.OutputFilm, error) {
	var films []domain.OutputFilm
	err := r.db.WithConnection(ctx, func(ctx context.Context, c *pgxpool.Conn) error {
		sqlStr := "SELECT DISTINCT films.id, films.title, films.description, films.release_date, films.rating FROM films "
		pagination := " LIMIT "

		var rows pgx.Rows
		var err error
		if filmTitleFragment != "" && actorNameFragment != "" {
			sqlStr += "JOIN film_actor ON films.id = film_actor.film_id JOIN actors ON actors.id = film_actor.actor_id WHERE films.title ILIKE '%' || $1 || '%' AND actors.name ILIKE '%' || $2 || '%'"
			pagination += "$3 OFFSET $4"
			rows, err = c.Query(ctx, sqlStr+pagination, filmTitleFragment, actorNameFragment, limit, (page-1)*limit)
		} else if actorNameFragment != "" {
			sqlStr += "JOIN film_actor ON films.id = film_actor.film_id JOIN actors ON actors.id = film_actor.actor_id WHERE actors.name ILIKE '%' || $1 || '%'"
			pagination += "$2 OFFSET $3"
			rows, err = c.Query(ctx, sqlStr+pagination, actorNameFragment, limit, (page-1)*limit)
		} else {
			sqlStr += "WHERE title ILIKE '%' || $1 || '%'"
			pagination += "$2 OFFSET $3"
			rows, err = c.Query(ctx, sqlStr+pagination, filmTitleFragment, limit, (page-1)*limit)
		}

		if err != nil {
			return err
		}

		var curFilm domain.OutputFilm
		var curReleaseDate time.Time
		for rows.Next() {
			err = rows.Scan(&curFilm.ID, &curFilm.Title, &curFilm.Description, &curReleaseDate, &curFilm.Rating)
			if err != nil {
				return err
			}

			curFilm.ReleaseDate = curReleaseDate.Format(time.DateOnly)

			actorRows, err := r.db.Query(ctx, "SELECT actor_id FROM film_actor WHERE film_id = $1 ORDER BY actor_id ASC", curFilm.ID)
			if err != nil {
				return err
			}

			var (
				curActor    domain.FilmOutputActor
				curGender   bool
				curBirthday time.Time
			)

			curFilm.Actors = make([]domain.FilmOutputActor, 0)
			for actorRows.Next() {
				err = actorRows.Scan(&curActor.ID)
				if err != nil {
					return err
				}

				err = r.db.QueryRow(ctx, "SELECT name, gender, birthday FROM actors WHERE id = $1", curActor.ID).Scan(&curActor.Name, &curGender, &curBirthday)
				if err != nil {
					return err
				}

				if curGender {
					curActor.Gender = "female"
				} else {
					curActor.Gender = "male"
				}

				curActor.Birthday = curBirthday.Format(time.DateOnly)

				curFilm.Actors = append(curFilm.Actors, curActor)
			}

			films = append(films, curFilm)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("repository.FindFilms(): %w", err)
	}

	return films, nil
}
