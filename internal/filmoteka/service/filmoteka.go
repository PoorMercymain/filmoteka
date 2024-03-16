package service

import (
	"context"
	"fmt"

	"github.com/PoorMercymain/filmoteka/internal/filmoteka/domain"
)

type filmoteka struct {
	repo domain.FilmotekaRepository
}

func New(repo domain.FilmotekaRepository) *filmoteka {
	return &filmoteka{repo: repo}
}

func (s *filmoteka) Ping(ctx context.Context) error {
	err := s.repo.Ping(ctx)
	if err != nil {
		return fmt.Errorf("service.Ping(): %w", err)
	}

	return nil
}
