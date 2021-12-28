package note_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/sksmith/note-server/core/note"
)

func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	mc := mockClock{}
	othertime, _ := time.Parse("2006-01-02", "2021-05-05")
	err := errors.New("some error")

	tests := []struct {
		ctx      context.Context
		input    note.Note
		repoErr  error
		repoNote note.Note
		wantErr  error
		wantNote note.Note
	}{
		{
			ctx:      context.Background(),
			input:    note.Note{ID: "id", Data: "some note"},
			wantNote: note.Note{ID: "id", Data: "some note", Created: mc.Now(), Updated: mc.Now()},
		},
		{
			ctx:      context.Background(),
			input:    note.Note{ID: "id", Data: "some note", Created: othertime, Updated: othertime},
			wantNote: note.Note{ID: "id", Data: "some note", Created: othertime, Updated: mc.Now()},
		},
		{
			ctx:     context.Background(),
			input:   note.Note{ID: "id", Data: "some note"},
			repoErr: err,
			wantErr: err,
		},
	}

	for _, test := range tests {
		mr := mockRepo{
			returnErr:  test.repoErr,
			returnNote: test.repoNote,
		}
		service := note.NewService(&mc, &mr)

		err := service.Create(test.ctx, test.input)
		if errors.Cause(err) != test.wantErr {
			t.Errorf("got=[%v] want=[%v]", err, test.wantErr)
		}
		if mr.savedNote != test.wantNote {
			t.Errorf("got=[%v] want=[%v]", mr.savedNote, test.wantNote)
		}
	}
}

func TestGet(t *testing.T) {
	mc := mockClock{}
	err := errors.New("some error")

	tests := []struct {
		ctx      context.Context
		input    string
		repoErr  error
		repoNote note.Note
		wantErr  error
		wantNote note.Note
	}{
		{
			ctx:      context.Background(),
			input:    "id",
			repoNote: note.Note{ID: "id", Data: "some note"},
			wantNote: note.Note{ID: "id", Data: "some note"},
		},
		{
			ctx:     context.Background(),
			input:   "id",
			repoErr: err,
			wantErr: err,
		},
	}

	for _, test := range tests {
		mr := mockRepo{
			returnErr:  test.repoErr,
			returnNote: test.repoNote,
		}
		service := note.NewService(&mc, &mr)

		got, err := service.Get(test.ctx, test.input)
		if errors.Cause(err) != test.wantErr {
			t.Errorf("got=[%v] want=[%v]", err, test.wantErr)
		}
		if got != test.wantNote {
			t.Errorf("got=[%v] want=[%v]", got, test.wantNote)
		}
	}
}

func TestDelete(t *testing.T) {
	mc := mockClock{}
	err := errors.New("some error")

	tests := []struct {
		ctx     context.Context
		input   string
		repoErr error
		wantErr error
	}{
		{
			ctx:   context.Background(),
			input: "id",
		},
		{
			ctx:     context.Background(),
			input:   "id",
			repoErr: err,
			wantErr: err,
		},
	}

	for _, test := range tests {
		mr := mockRepo{
			returnErr: test.repoErr,
		}
		service := note.NewService(&mc, &mr)

		err := service.Delete(test.ctx, test.input)
		if errors.Cause(err) != test.wantErr {
			t.Errorf("got=[%v] want=[%v]", err, test.wantErr)
		}
	}
}

func TestList(t *testing.T) {
	mc := mockClock{}
	err := errors.New("some error")

	tests := []struct {
		ctx           context.Context
		startIdx      int
		endIdx        int
		repoListNotes []note.ListNote
		repoErr       error
		wantListNotes []note.ListNote
		wantErr       error
	}{
		{
			ctx: context.Background(),
			repoListNotes: []note.ListNote{
				{ID: "1", Title: "Some Title"},
				{ID: "2", Title: "Some Other Title"},
			},
			wantListNotes: []note.ListNote{
				{ID: "1", Title: "Some Title"},
				{ID: "2", Title: "Some Other Title"},
			},
		},
		{
			ctx:     context.Background(),
			repoErr: err,
			wantErr: err,
		},
	}

	for _, test := range tests {
		mr := mockRepo{
			returnListNote: test.repoListNotes,
			returnErr:      test.repoErr,
		}
		service := note.NewService(&mc, &mr)

		got, err := service.List(test.ctx, test.startIdx, test.endIdx)
		if errors.Cause(err) != test.wantErr {
			t.Errorf("got=[%v] want=[%v]", err, test.wantErr)
		}
		if len(got) != len(test.wantListNotes) {
			t.Errorf("got=[%v] want=[%v]", len(got), len(test.wantListNotes))
		}
		for i, ln := range got {
			if test.wantListNotes[i] != ln {
				t.Errorf("got=[%v] want=[%v]", err, test.wantErr)
			}
		}
	}
}

type mockClock struct{}

func (m *mockClock) Now() time.Time {
	t, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z07:00")
	return t
}

type mockRepo struct {
	returnErr      error
	returnNote     note.Note
	returnListNote []note.ListNote
	savedNote      note.Note
}

func (r *mockRepo) Save(ctx context.Context, note note.Note) error {
	if r.returnErr != nil {
		return r.returnErr
	}
	r.savedNote = note
	return r.returnErr
}

func (r *mockRepo) Get(ctx context.Context, id string) (note.Note, error) {
	return r.returnNote, r.returnErr
}

func (r *mockRepo) Delete(ctx context.Context, id string) error {
	return r.returnErr
}

func (r *mockRepo) List(ctx context.Context, startIdx, endIdx int) ([]note.ListNote, error) {
	return r.returnListNote, r.returnErr
}
