package repository

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
)

func ApplyMigrations(filePath string, dsn string) error {
	m, err := migrate.New("file://"+filePath, dsn)
	if err != nil {
		return fmt.Errorf("repository.ApplyMigrations(): %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("repository.ApplyMigrations(): %w", err)
	}

	sourceErr, databaseErr := m.Close()

	if sourceErr != nil {
		return fmt.Errorf("repository.ApplyMigrations(): %w", sourceErr)
	}

	if databaseErr != nil {
		return fmt.Errorf("repository.ApplyMigrations(): %w", databaseErr)
	}

	return nil
}
