package noterepo

import (
	"context"

	"github.com/sksmith/note-server/core/note"
)

type MockRepo struct {
	SaveFunc   func(ctx context.Context, note *note.Note) error
	GetFunc    func(ctx context.Context, noteID string) (note.Note, error)
	UpdateFunc func(ctx context.Context, note *note.Note) error
	DeleteFunc func(ctx context.Context, noteID string) error
}

func (r MockRepo) Save(ctx context.Context, note *note.Note) error {
	return r.SaveFunc(ctx, note)
}

func (r MockRepo) Get(ctx context.Context, noteID string) (note.Note, error) {
	return r.GetFunc(ctx, noteID)
}

func (r MockRepo) Update(ctx context.Context, note *note.Note) error {
	return r.UpdateFunc(ctx, note)
}

func (r MockRepo) Delete(ctx context.Context, noteID string) error {
	return r.DeleteFunc(ctx, noteID)
}
