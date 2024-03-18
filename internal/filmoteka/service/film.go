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

func (s *film) CreateFilm(ctx context.Context, title string, description string, releaseDate time.Time, rating float32, actors []int) (int, error) {
	id, err := s.repo.CreateFilm(ctx, title, description, releaseDate, rating, actors)
	if err != nil {
		return 0, fmt.Errorf("service.CreateFilm(): %w", err)
	}

	return id, nil
}

func (s *film) UpdateFilm(ctx context.Context, id int, title string, description string, releaseDate time.Time, rating *float32, actors []int) error {
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

func (s *film) ReadFilms(ctx context.Context, field string, order string, page int, limit int) ([]domain.OutputFilm, error) {
	films, err := s.repo.ReadFilms(ctx, field, order, page, limit)
	if err != nil {
		return nil, fmt.Errorf("service.ReadFilms(): %w", err)
	}

	return films, nil
}

func (s *film) FindFilms(ctx context.Context, filmTitleFragment string, actorNameFragment string, page int, limit int) ([]domain.OutputFilm, error) {
	films, err := s.repo.FindFilms(ctx, filmTitleFragment, actorNameFragment, page, limit)
	if err != nil {
		return nil, fmt.Errorf("service.FindFilms(): %w", err)
	}

	return films, nil
}
