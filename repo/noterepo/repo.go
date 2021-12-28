package noterepo

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rs/zerolog/log"
	"github.com/sksmith/note-server/core"
	"github.com/sksmith/note-server/core/note"
)

type s3Repo struct {
	bucket     string
	uploader   Uploader
	downloader Downloader
	deleter    Deleter
}

const IndexID = "index"

type Downloader interface {
	Download(w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (n int64, err error)
}

type Uploader interface {
	Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
}

type Deleter interface {
	DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
}

func NewS3Repo(uploader Uploader, downloader Downloader, deleter Deleter, bucket string) *s3Repo {
	log.Info().
		Str("func", "NewS3Repo").
		Msg("setting up s3 session")

	return &s3Repo{
		bucket:     bucket,
		deleter:    deleter,
		uploader:   uploader,
		downloader: downloader,
	}
}

func (r *s3Repo) Save(ctx context.Context, note note.Note) error {
	n, err := json.Marshal(note)
	if err != nil {
		return err
	}
	_, err = r.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(note.ID),
		Body:   bytes.NewReader(n),
	})
	if err != nil {
		return err
	}

	// Note: There are no rollbacks with s3 storage so we can't rollback creating
	// the note if adding it to index fails.
	err = r.upsertNoteToIndex(ctx, note)
	if err != nil {
		return err
	}
	return nil
}

func (r *s3Repo) Get(ctx context.Context, id string) (note.Note, error) {
	data := aws.NewWriteAtBuffer([]byte{})
	s, err := r.downloader.Download(data, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(id),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return note.Note{}, &core.ErrNotFound{}
			default:
				return note.Note{}, err
			}
		} else {
			return note.Note{}, err
		}
	}

	log.Info().
		Str("func", "GetNote").
		Str("id", id).
		Int64("size", s).
		Msg("downloaded note note")

	n := note.Note{}
	err = json.Unmarshal(data.Bytes(), &n)
	if err != nil {
		return note.Note{}, err
	}

	return n, nil
}

func (r *s3Repo) Delete(ctx context.Context, id string) error {
	_, err := r.deleter.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(id),
	})

	if err != nil {
		return err
	}

	err = r.deleteNoteFromIndex(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *s3Repo) List(ctx context.Context, startIdx, endIdx int) ([]note.ListNote, error) {
	idx, err := r.getIndex(ctx)
	if err != nil {
		if core.IsErrNotFound(err) {
			return []note.ListNote{}, nil
		} else {
			return []note.ListNote{}, err
		}
	}

	_ = removeListNote(&idx, IndexID)

	return idx, nil
}

func (r *s3Repo) upsertNoteToIndex(ctx context.Context, n note.Note) error {
	list, err := r.List(ctx, 0, 0)
	if err != nil && !core.IsErrNotFound(err) {
		return err
	}

	updated := false
	for i := range list {
		if list[i].ID != n.ID {
			continue
		}

		list[i].Title = n.Title
		list[i].Created = n.Created
		list[i].Updated = n.Updated
		updated = true
		break
	}

	if !updated {
		list = append(list, mapNoteToListNote(n))
	}

	err = r.saveIndex(ctx, list)
	if err != nil {
		return err
	}

	return nil
}

func (r *s3Repo) deleteNoteFromIndex(ctx context.Context, ID string) error {
	list, err := r.List(ctx, 0, 0)
	if err != nil {
		return err
	}

	idx := removeListNote(&list, ID)
	if idx != -1 {
		err = r.saveIndex(ctx, list)
		if err != nil {
			return err
		}
	}

	return nil
}

// returns a -1 if ID not found
func removeListNote(l *[]note.ListNote, ID string) int {
	idx := -1
	list := *l
	for i, ln := range list {
		if ln.ID == ID {
			idx = i
			break
		}
	}

	if idx != -1 {
		*l = append(list[:idx], list[idx+1:]...)
	}

	return idx
}

func (r *s3Repo) getIndex(ctx context.Context) ([]note.ListNote, error) {
	data := aws.NewWriteAtBuffer([]byte{})
	s, err := r.downloader.Download(data, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(IndexID),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return []note.ListNote{}, &core.ErrNotFound{}
			default:
				return []note.ListNote{}, err
			}
		} else {
			return []note.ListNote{}, err
		}
	}

	log.Info().
		Str("func", "getIndex").
		Int64("size", s).
		Msg("downloaded index")

	l := make([]note.ListNote, 0)
	err = json.Unmarshal(data.Bytes(), &l)
	if err != nil {
		return []note.ListNote{}, err
	}

	return l, nil
}

func (r *s3Repo) saveIndex(ctx context.Context, idx []note.ListNote) error {
	data, err := json.Marshal(idx)
	if err != nil {
		return err
	}
	_, err = r.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(IndexID),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return err
	}

	return nil
}

func mapNoteToListNote(n note.Note) note.ListNote {
	return note.ListNote{
		ID:      n.ID,
		Title:   n.Title,
		Created: n.Created,
		Updated: n.Updated,
	}
}
