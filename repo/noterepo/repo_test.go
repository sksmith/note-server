package noterepo_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rs/zerolog"
	"github.com/sksmith/note-server/core"
	"github.com/sksmith/note-server/core/note"
	"github.com/sksmith/note-server/repo/noterepo"
)

func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	os.Exit(m.Run())
}

func TestGet(t *testing.T) {
	err := errors.New("some unkown error")
	tests := []struct {
		name     string
		ctx      context.Context
		input    string
		s3Err    error
		s3Note   string
		wantNote note.Note
		wantErr  error
	}{
		{
			name:     "Happy Path",
			input:    "1",
			s3Note:   `{"id": "1", "data": "somenote"}`,
			wantNote: note.Note{ID: "1", Data: "somenote"},
		},
		{
			name:  "Key Not Found",
			input: "1",
			s3Err: awserr.New(
				s3.ErrCodeNoSuchKey,
				"no such key",
				errors.New("madeup error"),
			),
			wantErr: &core.ErrNotFound{},
		},
		{
			name:    "Unknown Error",
			input:   "1",
			s3Err:   err,
			wantErr: err,
		},
	}

	for _, test := range tests {
		downloader := &mockDownloader{note: test.s3Note, err: test.s3Err}
		uploader := &mockUploader{}
		deleter := &mockDeleter{}

		repo := noterepo.NewS3Repo(uploader, downloader, deleter, "somebucket")
		got, err := repo.Get(test.ctx, test.input)

		compare(test.name, err, test.wantErr, t)
		compare(test.name, got, test.wantNote, t)
	}
}

func TestList(t *testing.T) {
	err := errors.New("some unkown error")

	tests := []struct {
		name          string
		ctx           context.Context
		startIdx      int
		endIdx        int
		s3Err         error
		s3ListNotes   string
		wantListNotes []note.ListNote
		wantErr       error
	}{
		{
			name:        "Happy Path",
			s3ListNotes: `[{"id": "1", "title": "somenote"},{"id": "2", "title": "someothernote"}]`,
			wantListNotes: []note.ListNote{
				{ID: "1", Title: "somenote"},
				{ID: "2", Title: "someothernote"},
			},
		},
		{
			name: "Not Found",
			s3Err: awserr.New(
				s3.ErrCodeNoSuchKey,
				"no such key",
				errors.New("madeup error"),
			),
			wantErr:       nil,
			wantListNotes: []note.ListNote{},
		},
		{
			name:    "Unknown Error",
			s3Err:   err,
			wantErr: err,
		},
	}

	for _, test := range tests {
		downloader := &mockDownloader{listNotes: test.s3ListNotes, err: test.s3Err}
		uploader := &mockUploader{}
		deleter := &mockDeleter{}
		repo := noterepo.NewS3Repo(uploader, downloader, deleter, "somebucket")

		got, err := repo.List(test.ctx, test.startIdx, test.endIdx)

		compare(test.name, err, test.wantErr, t)
		compare(test.name, len(got), len(test.wantListNotes), t)

		for i, ln := range got {
			compare(test.name, test.wantListNotes[i], ln, t)
		}
	}
}

func TestSave(t *testing.T) {
	err := errors.New("some unkown error")

	tests := []struct {
		name          string
		ctx           context.Context
		input         note.Note
		s3Err         error
		s3Note        string
		s3ListNotes   string
		wantNote      string
		wantListNotes string
		wantErr       error
	}{
		{
			name:          "Add a New Note",
			input:         note.Note{ID: "2", Title: "some note title", Data: "some other note"},
			s3ListNotes:   marshal([]note.ListNote{{ID: "1", Title: "somenote"}}),
			wantNote:      marshal(note.Note{ID: "2", Title: "some note title", Data: "some other note"}),
			wantListNotes: marshal([]note.ListNote{{ID: "1", Title: "somenote"}, {ID: "2", Title: "some note title"}}),
		},
		{
			name:          "Update a Note",
			input:         note.Note{ID: "1", Title: "some updated note", Data: "some new text"},
			s3ListNotes:   marshal([]note.ListNote{{ID: "1", Title: "somenote"}, {ID: "2", Title: "some other note"}}),
			wantNote:      marshal(note.Note{ID: "1", Title: "some updated note", Data: "some new text"}),
			wantListNotes: marshal([]note.ListNote{{ID: "1", Title: "some updated note"}, {ID: "2", Title: "some other note"}}),
		},
		{
			name:    "Unknown Error",
			input:   note.Note{},
			s3Err:   err,
			wantErr: err,
		},
	}

	for _, test := range tests {
		downloader := &mockDownloader{note: test.s3Note, listNotes: test.s3ListNotes, err: test.s3Err}
		uploader := &mockUploader{err: test.s3Err}
		deleter := &mockDeleter{}

		repo := noterepo.NewS3Repo(uploader, downloader, deleter, "somebucket")
		err := repo.Save(test.ctx, test.input)

		compare(test.name, err, test.wantErr, t)
		compare(test.name, uploader.uploadedNote, test.wantNote, t)
		compare(test.name, uploader.uploadedIndex, test.wantListNotes, t)
	}
}

func TestDelete(t *testing.T) {
	err := errors.New("some unkown error")

	tests := []struct {
		name          string
		ctx           context.Context
		input         string
		s3Err         error
		s3ListNotes   string
		wantListNotes string
		wantErr       error
	}{
		{
			name:          "Delete a Note",
			input:         "1",
			s3ListNotes:   marshal([]note.ListNote{{ID: "1", Title: "somenote"}, {ID: "2", Title: "some note title"}}),
			wantListNotes: marshal([]note.ListNote{{ID: "2", Title: "some note title"}}),
		},
		{
			name:          "Delete a Missing Note",
			input:         "3",
			s3ListNotes:   marshal([]note.ListNote{{ID: "1", Title: "somenote"}, {ID: "2", Title: "some note title"}}),
			wantListNotes: "",
		},
		{
			name:    "Unknown Error",
			input:   "1",
			s3Err:   err,
			wantErr: err,
		},
	}

	for _, test := range tests {
		downloader := &mockDownloader{listNotes: test.s3ListNotes, err: test.s3Err}
		uploader := &mockUploader{err: test.s3Err}
		deleter := &mockDeleter{err: test.s3Err}

		repo := noterepo.NewS3Repo(uploader, downloader, deleter, "somebucket")
		err := repo.Delete(test.ctx, test.input)

		compare(test.name, err, test.wantErr, t)
		compare(test.name, deleter.requestedID, test.input, t)
		compare(test.name, uploader.uploadedIndex, test.wantListNotes, t)
	}
}

func compare(testName string, got, want interface{}, t *testing.T) {
	if got != want {
		t.Errorf("%v: got=[%v] want=[%v]", testName, got, want)
	}
}

func marshal(v interface{}) string {
	val, _ := json.Marshal(v)
	return string(val)
}

type mockDownloader struct {
	note      string
	listNotes string
	err       error
}

func (m *mockDownloader) Download(w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (n int64, err error) {
	if m.err != nil {
		return -1, m.err
	}
	if *input.Key == noterepo.IndexID {
		_, _ = w.WriteAt([]byte(m.listNotes), 0)
	} else {
		_, _ = w.WriteAt([]byte(m.note), 0)
	}
	return -1, nil
}

type mockUploader struct {
	uploadedNote  string
	uploadedIndex string
	err           error
}

func (m *mockUploader) Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	if m.err != nil {
		return &s3manager.UploadOutput{}, m.err
	}
	if *input.Key == noterepo.IndexID {
		bytes, err := ioutil.ReadAll(input.Body)
		if err != nil {
			return &s3manager.UploadOutput{}, err
		}
		m.uploadedIndex = string(bytes)
	} else {
		bytes, err := ioutil.ReadAll(input.Body)
		if err != nil {
			return &s3manager.UploadOutput{}, err
		}
		m.uploadedNote = string(bytes)
	}
	return &s3manager.UploadOutput{}, nil
}

type mockDeleter struct {
	requestedID string
	err         error
}

func (m *mockDeleter) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	m.requestedID = *input.Key
	if m.err != nil {
		return &s3.DeleteObjectOutput{}, m.err
	}
	return &s3.DeleteObjectOutput{}, nil
}
