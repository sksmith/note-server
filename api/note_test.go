package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
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

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("expected %v got %v", http.StatusOK, w.Result().StatusCode)
	}
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

func TestDeleteNotFound(t *testing.T) {
	r := httptest.NewRequest(http.MethodDelete, "/1", nil)
	w := httptest.NewRecorder()

	svc := mockNoteService{}
	svc.ReturnError(&core.ErrNotFound{})
	noteApi := api.NewNoteApi(svc)

	noteApi.Delete(w, r)

	if w.Result().StatusCode != http.StatusNoContent {
		t.Errorf("expected %v got %v", http.StatusNoContent, w.Result().StatusCode)
	}
}

func TestDeleteInternalServerError(t *testing.T) {
	r := httptest.NewRequest(http.MethodDelete, "/1", nil)
	w := httptest.NewRecorder()

	svc := mockNoteService{}
	svc.ReturnError(errors.New("some unexpected error"))
	noteApi := api.NewNoteApi(svc)

	noteApi.Delete(w, r)

	if w.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf("expected %v got %v", http.StatusInternalServerError, w.Result().StatusCode)
	}
}

func TestList(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	noteApi := api.NewNoteApi(mockNoteService{})

	noteApi.List(w, r)
	lr := parseListResponse(w, t)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("expected %v got %v", http.StatusOK, w.Result().StatusCode)
	}
	if len(lr.Notes) != 2 {
		t.Errorf("expected %v got %v", 2, len(lr.Notes))
	}
}

func TestListInternalServerError(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	svc := mockNoteService{}
	svc.ReturnError(errors.New("some unexpected error"))
	noteApi := api.NewNoteApi(svc)

	noteApi.List(w, r)
	_ = parseErrorResponse(w, t)

	if w.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf("expected %v got %v", http.StatusInternalServerError, w.Result().StatusCode)
	}
}

func TestCreate(t *testing.T) {
	r := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{"id": "1", "data": "somenote"}`))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	noteApi := api.NewNoteApi(mockNoteService{})

	noteApi.Create(w, r)
	n := parseResponse(w, t)

	if w.Result().StatusCode != http.StatusCreated {
		t.Errorf("expected %v got %v", http.StatusCreated, w.Result().StatusCode)
	}
	if n.ID != "1" {
		t.Errorf("expected %v got %v", "1", n.ID)
	}
}

func TestCreateBadRequest(t *testing.T) {
	r := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{"id": "1", "badfield": "somenote"}`))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	noteApi := api.NewNoteApi(mockNoteService{})

	noteApi.Create(w, r)
	_ = parseErrorResponse(w, t)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected %v got %v", http.StatusBadRequest, w.Result().StatusCode)
	}
}

func TestCreateInternalServerError(t *testing.T) {
	r := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{"id": "1", "data": "somenote"}`))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	svc := mockNoteService{}
	svc.ReturnError(errors.New("some unexpected exception"))
	noteApi := api.NewNoteApi(svc)

	noteApi.Create(w, r)
	_ = parseErrorResponse(w, t)

	if w.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf("expected %v got %v", http.StatusInternalServerError, w.Result().StatusCode)
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
		{ID: "2"},
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

func parseResponse(w *httptest.ResponseRecorder, t *testing.T) api.NoteResponse {
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	n := api.NoteResponse{}
	err = json.Unmarshal(data, &n)
	if err != nil {
		t.Errorf("failed to parse response %v", err)
	}
	return n
}

func parseListResponse(w *httptest.ResponseRecorder, t *testing.T) api.ListNoteResponse {
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	ln := api.ListNoteResponse{}
	err = json.Unmarshal(data, &ln)
	if err != nil {
		t.Errorf("failed to parse response %v", err)
	}
	return ln
}
