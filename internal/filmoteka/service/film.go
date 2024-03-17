package service

import (
	"context"
	"fmt"
	"time"

	"github.com/PoorMercymain/filmoteka/internal/filmoteka/domain"
)

var (
	_ domain.FilmService = (*film)(nil)
)

type film struct {
	repo domain.FilmRepository
}

func NewFilm(repo domain.FilmRepository) *film {
	return &film{repo: repo}
}

func (s *film) Ping(ctx context.Context) error {
	err := s.repo.Ping(ctx)
	if err != nil {
		return fmt.Errorf("service.Ping(): %w", err)
	}

	return nil
}

func (s *film) CreateFilm(ctx context.Context, title string, description string, releaseDate time.Time, rating int, actors []int) (int, error) {
	id, err := s.repo.CreateFilm(ctx, title, description, releaseDate, rating, actors)
	if err != nil {
		return 0, fmt.Errorf("service.CreateFilm(): %w", err)
	}

	return id, nil
}
