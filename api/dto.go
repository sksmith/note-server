package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/sksmith/note-server/config"
	"github.com/sksmith/note-server/core/note"
)

type ListNoteResponse struct {
	Notes []note.ListNote `json:"notes"`
}

func NewListNoteResponse(l []note.ListNote) *ListNoteResponse {
	resp := &ListNoteResponse{Notes: l}
	return resp
}

func (nr *ListNoteResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
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

type CreateNoteRequest struct {
	*note.Note
}

func (p *CreateNoteRequest) Bind(_ *http.Request) error {
	if p.Note.ID == "" || p.Note.Data == "" {
		return errors.New("missing required field(s)")
	}

	return nil
}

func Render(w http.ResponseWriter, r *http.Request, rnd render.Renderer) {
	if err := render.Render(w, r, rnd); err != nil {
		log.Warn().Err(err).Msg("failed to render")
	}
}

type EnvResponse struct {
	config.Config
}

func NewEnvResponse(c config.Config) *EnvResponse {
	resp := &EnvResponse{Config: c}
	return resp
}

func (er *EnvResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}
