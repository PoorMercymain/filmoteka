package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	appErrors "github.com/PoorMercymain/filmoteka/errors"
	"github.com/PoorMercymain/filmoteka/internal/filmoteka/domain"
)

var (
	_ domain.ActorRepository = (*actor)(nil)
)

type actor struct {
	db *postgres
}

func NewActor(pg *postgres) *actor {
	return &actor{db: pg}
}

func (r *actor) CreateActor(ctx context.Context, name string, gender bool, birthday time.Time) (int, error) {
	var id int
	err := r.db.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		return tx.QueryRow(ctx, "INSERT INTO actors(name, gender, birthday) VALUES($1, $2, $3) RETURNING id", name, gender, birthday).Scan(&id)
	})

	if err != nil {
		return 0, fmt.Errorf("repository.CreateActor(): %w", err)
	}

	return id, nil
}

func (r *actor) UpdateActor(ctx context.Context, id int, name string, gender *bool, birthday time.Time) error {
	err := r.db.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		var (
			nameInDB     string
			genderInDB   bool
			birthdayInDB time.Time
		)

		err := tx.QueryRow(ctx, "SELECT name, gender, birthday FROM actors WHERE id = $1", id).Scan(&nameInDB, &genderInDB, &birthdayInDB)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return appErrors.ErrNotFoundInDB
			}

			return err
		}

		if name == "" {
			name = nameInDB
		}

		if gender == nil {
			*gender = genderInDB
		}

		var timeDefault time.Time
		if birthday == timeDefault {
			birthday = birthdayInDB
		}

		_, err = tx.Exec(ctx, "UPDATE actors SET name = $1, gender = $2, birthday = $3 WHERE id = $4", name, *gender, birthday, id)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("repository.UpdateActor(): %w", err)
	}

	return nil
}

func (r *actor) DeleteActor(ctx context.Context, id int) error {
	err := r.db.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, "DELETE FROM actors WHERE id = $1", id)
		if err != nil {
			return err
		}

		if tag.RowsAffected() == 0 {
			return appErrors.ErrNotFoundInDB
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("repository.DeleteActor(): %w", err)
	}

	return nil
}
