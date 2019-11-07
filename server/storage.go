package main

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/s3blob"
	"io"
)

type DocumentSave interface {
	Save(ctx context.Context, fileName string, reader io.Reader) (path string, err error)
}

type DocumentGet interface {
	Get(ctx context.Context, writer io.Writer) error
}

type DocumentStorage interface {
	DocumentSave
	DocumentGet
}

type BucketStorage struct {
	Bucket *blob.Bucket
}

func NewBucketStorage(config BucketConfig) *BucketStorage {
	bucket, err := blob.OpenBucket(context.Background(), config.ConnectionString)
	if err != nil {
		logrus.WithError(err).Fatal("unable to connect to bucket")
	}
	return &BucketStorage{
		Bucket: bucket,
	}
}

func NewBucketDocumentStorage(storage *BucketStorage) DocumentSave {
	return storage
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
		return "", errors.Wrap(err, "failed to upload file")
	}

	return "", err
}
