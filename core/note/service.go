package note

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

type service struct {
	repo Repository
}

func (s *service) Create(ctx context.Context, note Note) error {
	const funcName = "CreateNote"

	log.Info().
		Str("func", funcName).
		Str("id", note.ID).
		Msg("creating note")

	if note.Created.IsZero() {
		note.Created = time.Now().UTC()
	}
	note.Updated = time.Now().UTC()

	if err := s.repo.Save(ctx, note); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *service) Get(ctx context.Context, id string) (Note, error) {
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

func (s *service) Delete(ctx context.Context, id string) error {
	const funcName = "DeleteNote"

	log.Info().
		Str("func", funcName).
		Str("id", id).
		Msg("deleting note")

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *service) List(ctx context.Context, startIdx, endIdx int) ([]ListNote, error) {
	const funcName = "ListNote"

	log.Info().
		Str("func", funcName).
		Msg("listing notes")

	list, err := s.repo.List(ctx, startIdx, endIdx)
	if err != nil {
		return []ListNote{}, errors.WithStack(err)
	}

	return list, nil
}

type Repository interface {
	Save(ctx context.Context, note Note) error
	Get(ctx context.Context, id string) (Note, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, startIdx, endIdx int) ([]ListNote, error)
}
