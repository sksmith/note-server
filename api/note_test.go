package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sksmith/note-server/api"
	"github.com/sksmith/note-server/core/note"
)

func TestGet(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/1", nil)
	w := httptest.NewRecorder()

	noteApi := api.NewNoteApi(mockNoteService{})

	noteApi.Get(w, r)
	resp := parseResponse(w, t)

	if resp.ID != "1" {
		t.Errorf("expected 1 got %v", resp.ID)
	}
}

func TestGetError(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/upper?word=abc", nil)
	w := httptest.NewRecorder()

	noteApi := api.NewNoteApi(mockErrorNoteService{})

	noteApi.Get(w, r)
	errResp := parseErrorResponse(w, t)

	if w.Result().StatusCode != 500 {
		t.Errorf("expected 500 got %v", w.Result().StatusCode)
	}
	if errResp.ErrorText != "An internal server error has occurred." {
		t.Errorf("expected \"An internal server error has occurred.\" got %v", errResp.ErrorText)
	}
}

func TestDelete(t *testing.T) {
	r := httptest.NewRequest(http.MethodDelete, "/1", nil)
	w := httptest.NewRecorder()

	noteApi := api.NewNoteApi(mockNoteService{})

	noteApi.Delete(w, r)

	if w.Result().StatusCode != 200 {
		t.Errorf("expected 500 got %v", w.Result().StatusCode)
	}
}

type mockNoteService struct{}

func (m mockNoteService) GetNote(context.Context, string) (note.Note, error) {
	return note.Note{
		ID:   "1",
		Note: "somenote",
	}, nil
}

func (m mockNoteService) CreateNote(context.Context, note.Note) error {
	return nil
}

func (m mockNoteService) DeleteNote(context.Context, string) error {
	return nil
}

type mockErrorNoteService struct{}

func (m mockErrorNoteService) GetNote(_ context.Context, _ string) (note.Note, error) {
	return note.Note{}, errors.New("something horrible went wrong")
}

func (m mockErrorNoteService) CreateNote(_ context.Context, _ note.Note) error {
	return errors.New("something horrible happened")
}

func (m mockErrorNoteService) DeleteNote(_ context.Context, _ string) error {
	return errors.New("made up delete failure")
}

func parseErrorResponse(w *httptest.ResponseRecorder, t *testing.T) api.ErrResponse {
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	e := api.ErrResponse{}
	err = json.Unmarshal(data, &e)
	if err != nil {
		t.Errorf("failed to parse response %v", err)
	}
	return e
}

func parseResponse(w *httptest.ResponseRecorder, t *testing.T) note.Note {
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	n := note.Note{}
	err = json.Unmarshal(data, &n)
	if err != nil {
		t.Errorf("failed to parse response %v", err)
	}
	return n
}
