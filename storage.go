package main

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gocloud.dev/blob"
	"io"
)

type BookSave interface {
	Save(ctx context.Context, reader io.Reader) (string, error)
}

type BookGet interface {
	Get(ctx context.Context, writer io.Writer) error
}

type BookStorage interface {
	BookSave
	BookGet
}

type BucketStorage struct {
	Bucket blob.Bucket
}

func NewBucketStorage(bucket blob.Bucket) BucketStorage {
	return BucketStorage{
		Bucket:bucket,
	}
}

func (s *BucketStorage) Save(ctx context.Context, fileName string, reader io.Reader) (path string, err error) {

	w, err := s.Bucket.NewWriter(ctx, fileName, nil)
	if err != nil {
		logrus.WithError(err).Error("unable to create upload writer")
		return "", errors.Wrap(err, "unable to creat upload writer")
	}

	defer w.Close()

	if _, err := io.Copy(w, reader); err != nil {
		logrus.WithError(err).Error("failed to upload file")
		return "",errors.Wrap(err, "failed to upload file")
	}

	return "", err
}
