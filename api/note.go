package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sksmith/note-server/core"
	"github.com/sksmith/note-server/core/note"
)

type NoteApi struct {
	service NoteService
}

type NoteService interface {
	Get(context.Context, string) (note.Note, error)
	Create(context.Context, note.Note) error
	Delete(context.Context, string) error
	List(context.Context, int, int) ([]note.ListNote, error)
}

func NewNoteApi(service NoteService) *NoteApi {
	return &NoteApi{service: service}
}

const (
	CtxKeyProduct     CtxKey = "product"
	CtxKeyReservation CtxKey = "reservation"
)

func (n *NoteApi) ConfigureRouter(r chi.Router) {
	r.Get("/", n.List)
	r.Put("/", n.Create)
	r.Get("/{id}", n.Get)
	r.Delete("/{id}", n.Delete)
}

func (a *NoteApi) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	n, err := a.service.Get(r.Context(), id)
	if err != nil {
		handleError(w, r, err)
		return
	}

	Render(w, r, NewNoteResponse(n))
}

func (a *NoteApi) List(w http.ResponseWriter, r *http.Request) {
	n, err := a.service.List(r.Context(), 0, 0)
	if err != nil {
		handleError(w, r, err)
		return
	}

	Render(w, r, NewListNoteResponse(n))
}

func (a *NoteApi) Create(w http.ResponseWriter, r *http.Request) {
	data := &CreateNoteRequest{}
	if err := render.Bind(r, data); err != nil {
		Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := a.service.Create(r.Context(), *data.Note); err != nil {
		handleError(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	Render(w, r, NewNoteResponse(*data.Note))
}

func (a *NoteApi) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := a.service.Delete(r.Context(), id)
	if err != nil {
		handleError(w, r, err)
		return
	}

	render.NoContent(w, r)
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	log.Err(err).Send()
	switch errors.Cause(err).(type) {
	case *core.ErrNotFound:
		Render(w, r, ErrNotFound)
	default:
		Render(w, r, ErrInternalServer)
	}
}
