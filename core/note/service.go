package note

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

type Service interface {
	GetNote(ctx context.Context, sku string) (Note, error)
	CreateNote(ctx context.Context, product Note) error
}

type service struct {
	repo Repository
}

func (s *service) CreateNote(ctx context.Context, note Note) error {
	const funcName = "CreateNote"

	log.Info().
		Str("func", funcName).
		Str("id", note.ID).
		Msg("creating note")

	if err := s.repo.Save(ctx, note); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *service) GetNote(ctx context.Context, id string) (Note, error) {
	const funcName = "GetNote"

	log.Info().
		Str("func", funcName).
		Str("id", id).
		Msg("getting note")

	note, err := s.repo.Get(ctx, id)
	if err != nil {
		return note, errors.WithStack(err)
	}
	return note, nil
}

type Repository interface {
	Save(ctx context.Context, note Note) error
	Get(ctx context.Context, sku string) (Note, error)
}
