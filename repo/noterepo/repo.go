package noterepo

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rs/zerolog/log"
	"github.com/sksmith/note-server/core/note"
)

type s3Repo struct {
	bucket     string
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

func New() note.Repository {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))
	return &s3Repo{
		bucket:     "",
		uploader:   s3manager.NewUploader(sess),
		downloader: s3manager.NewDownloader(sess),
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
	return nil
}

func (r *s3Repo) Get(ctx context.Context, id string) (note.Note, error) {
	data := aws.NewWriteAtBuffer([]byte{})
	n, err := r.downloader.Download(data, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(id),
	})
	if err != nil {
		return note.Note{}, err
	}

	log.Info().
		Str("func", "GetNote").
		Str("id", id).
		Str("note", string(data.Bytes())).
		Int64("size", n).
		Msg("downloaded note note")

	note := note.Note{}
	json.Unmarshal(data.Bytes(), &note)

	return note, nil
}

func (r *s3Repo) Delete(ctx context.Context, id string) error {

	return nil
}
