package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/sksmith/note-server/core/note"
)

type NoteApi struct {
	service note.Service
}

func NewNoteApi(service note.Service) *NoteApi {
	return &NoteApi{service: service}
}

const (
	CtxKeyProduct     CtxKey = "product"
	CtxKeyReservation CtxKey = "reservation"
)

func (n *NoteApi) ConfigureRouter(r chi.Router) {
	r.Put("/", n.Create)
	r.Get("/{id}", n.Get)
}

type NoteResponse struct {
	note.Note
}

func NewNoteResponse(n note.Note) *NoteResponse {
	resp := &NoteResponse{Note: n}
	return resp
}

func (nr *NoteResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

func (a *NoteApi) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	n, err := a.service.GetNote(r.Context(), id)
	if err != nil {
		log.Err(err).Send()
		Render(w, r, ErrInternalServer)
		return
	}

	render.Status(r, http.StatusOK)
	Render(w, r, NewNoteResponse(n))
}

func (a *NoteApi) Create(w http.ResponseWriter, r *http.Request) {
	data := &CreateNoteRequest{}
	if err := render.Bind(r, data); err != nil {
		Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := a.service.CreateNote(r.Context(), *data.Note); err != nil {
		log.Err(err).Send()
		Render(w, r, ErrInternalServer)
		return
	}

	render.Status(r, http.StatusCreated)
	Render(w, r, NewNoteResponse(*data.Note))
}

type CreateNoteRequest struct {
	*note.Note
}

func (p *CreateNoteRequest) Bind(_ *http.Request) error {
	if p.Note.ID == "" || p.Note.Note == "" {
		return errors.New("missing required field(s)")
	}

	return nil
}

func Render(w http.ResponseWriter, r *http.Request, rnd render.Renderer) {
	if err := render.Render(w, r, rnd); err != nil {
		log.Warn().Err(err).Msg("failed to render")
	}
}

func RenderList(w http.ResponseWriter, r *http.Request, l []render.Renderer) {
	if err := render.RenderList(w, r, l); err != nil {
		log.Warn().Err(err).Msg("failed to render")
	}
}
