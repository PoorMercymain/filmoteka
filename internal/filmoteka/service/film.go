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

func (s *film) UpdateFilm(ctx context.Context, id int, title string, description string, releaseDate time.Time, rating *int, actors []int) error {
	err := s.repo.UpdateFilm(ctx, id, title, description, releaseDate, rating, actors)
	if err != nil {
		return fmt.Errorf("service.UpdateFilm(): %w", err)
	}

	return nil
}

func (s *film) DeleteFilm(ctx context.Context, id int) error {
	err := s.repo.DeleteFilm(ctx, id)
	if err != nil {
		return fmt.Errorf("service.DeleteFilm(): %w", err)
	}

	return nil
}
