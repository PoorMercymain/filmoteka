package service

import (
	"context"
	"fmt"
	"time"

	"github.com/PoorMercymain/filmoteka/internal/filmoteka/domain"
)

var (
	_ domain.ActorService = (*actor)(nil)
)

type actor struct {
	repo domain.ActorRepository
}

func NewActor(repo domain.ActorRepository) *actor {
	return &actor{repo: repo}
}

func (s *actor) CreateActor(ctx context.Context, name string, gender bool, birthday time.Time) (int, error) {
	id, err := s.repo.CreateActor(ctx, name, gender, birthday)
	if err != nil {
		return 0, fmt.Errorf("service.CreateActor(): %w", err)
	}

	return id, nil
}

func (s *actor) UpdateActor(ctx context.Context, id int, name string, gender *bool, birthday time.Time) error {
	err := s.repo.UpdateActor(ctx, id, name, gender, birthday)
	if err != nil {
		return fmt.Errorf("service.UpdateActor(): %w", err)
	}

	return nil
}

func (s *actor) DeleteActor(ctx context.Context, id int) error {
	err := s.repo.DeleteActor(ctx, id)
	if err != nil {
		return fmt.Errorf("service.DeleteActor(): %w", err)
	}

	return nil
}
