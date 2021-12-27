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
	"github.com/sksmith/note-server/core"
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

func TestGetInternalServerError(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/1", nil)
	w := httptest.NewRecorder()

	svc := mockNoteService{}
	svc.ReturnError(errors.New("some weird error"))
	noteApi := api.NewNoteApi(svc)

	noteApi.Get(w, r)
	errResp := parseErrorResponse(w, t)

	if w.Result().StatusCode != 500 {
		t.Errorf("expected 500 got %v", w.Result().StatusCode)
	}
	if errResp.ErrorText != "An internal server error has occurred." {
		t.Errorf("expected \"An internal server error has occurred.\" got %v", errResp.ErrorText)
	}
}

func TestGetNotFoundError(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/1", nil)
	w := httptest.NewRecorder()

	svc := mockNoteService{}
	svc.ReturnError(&core.ErrNotFound{})
	noteApi := api.NewNoteApi(svc)

	noteApi.Get(w, r)

	if w.Result().StatusCode != 404 {
		t.Errorf("expected 404 got %v", w.Result().StatusCode)
	}
}

func TestDelete(t *testing.T) {
	r := httptest.NewRequest(http.MethodDelete, "/1", nil)
	w := httptest.NewRecorder()

	noteApi := api.NewNoteApi(mockNoteService{})

	noteApi.Delete(w, r)

	if w.Result().StatusCode != 204 {
		t.Errorf("expected 204 got %v", w.Result().StatusCode)
	}
}

type mockNoteService struct {
	returnError error
}

func (m *mockNoteService) ReturnError(err error) {
	m.returnError = err
}

func (m mockNoteService) Get(context.Context, string) (note.Note, error) {
	if m.returnError != nil {
		return note.Note{}, m.returnError
	}
	return note.Note{
		ID:   "1",
		Data: "somenote",
	}, nil
}

func (m mockNoteService) Create(context.Context, note.Note) error {
	if m.returnError != nil {
		return m.returnError
	}
	return nil
}

func (m mockNoteService) Delete(context.Context, string) error {
	if m.returnError != nil {
		return m.returnError
	}
	return nil
}

func (m mockNoteService) List(ctx context.Context, startIdx, endIdx int) ([]note.ListNote, error) {
	if m.returnError != nil {
		return []note.ListNote{}, m.returnError
	}
	return []note.ListNote{
		{ID: "1"},
	}, nil
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
