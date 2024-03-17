package service

import (
	"context"
	"fmt"
	"time"

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

func (s *filmoteka) CreateActor(ctx context.Context, name string, gender bool, birthday time.Time) (int, error) {
	id, err := s.repo.CreateActor(ctx, name, gender, birthday)
	if err != nil {
		return 0, fmt.Errorf("service.CreateActor(): %w", err)
	}

	return id, nil
}
