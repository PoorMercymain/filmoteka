package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type postgres struct {
	*pgxpool.Pool
}

func NewPostgres(pool *pgxpool.Pool) *postgres {
	return &postgres{pool}
}

func GetPgxPool(DSN string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(DSN)
	if err != nil {
		return nil, fmt.Errorf("repository.GetPgxPool(): %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("repository.GetPgxPool(): %w", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("repository.GetPgxPool(): %w", err)
	}

	return pool, nil
}

func (p *postgres) WithTransaction(ctx context.Context, txFunc func(context.Context, pgx.Tx) error) error {
	conn, err := p.Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("repository.WithTransaction(): %w", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("repository.WithTransaction(): %w", err)
	}

	err = txFunc(ctx, tx)
	if err != nil {
		rollbackErr := p.Rollback(ctx, tx)
		if rollbackErr != nil {
			resErr := errors.Join(err, rollbackErr)
			return fmt.Errorf("repository.WithTransaction(): %w", resErr)
		}

		return fmt.Errorf("repository.WithTransaction(): %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("repository.WithTransaction(): %w", err)
	}

	return nil
}

func (p *postgres) Rollback(ctx context.Context, tx pgx.Tx) error {
	err := tx.Rollback(ctx)
	if !errors.Is(err, pgx.ErrTxClosed) && err != nil {
		return fmt.Errorf("repository.Rollback(): %w", err)
	}

	return nil
}

func (p *postgres) WithConnection(ctx context.Context, connFunc func(context.Context, *pgxpool.Conn) error) error {
	conn, err := p.Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("repository.WithConnection(): %w", err)
	}
	defer conn.Release()

	err = connFunc(ctx, conn)
	if err != nil {
		return fmt.Errorf("repository.WithConnection(): %w", err)
	}

	return nil
}
