package main

import (
	"context"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gocloud.dev/blob"
	"io"
	"mime/multipart"
)

type BookUploader interface {
	Upload(ctx context.Context, reader io.Reader)
}

type BookDownloader interface {
	Download(ctx context.Context, writer io.Writer)
}

type BucketStorage struct {
	Bucket blob.Bucket
}

func NewBucketStorage(bucket blob.Bucket) BucketStorage {
	return BucketStorage{
		Bucket:bucket,
	}
}

func (s *BucketStorage) Upload(ctx context.Context, fileName string, reader multipart.File) (err error) {

	head := make([]byte, 261)
	if bytesRead, err := io.ReadFull(reader, head); err == io.EOF {
		logrus.WithField("bytesRead", bytesRead).WithError(err).Error("couldn't read file header: unexpected EOF")
		return io.ErrUnexpectedEOF
	} else if err != nil {
		logrus.WithField("bytesRead", bytesRead).WithError(err).Error("couldn't read file header")
		return err
	}

	reader.Seek(0, io.SeekStart)


	if kind, err := filetype.Match(head); err != nil {
		logrus.WithError(err).Error("unable to determine file type")
		return errors.Wrap(err, "unable to determine file type")
	} else {
		if kind != matchers.TypeEpub && kind != matchers.TypePdf {
			logrus.WithFields(logrus.Fields{"mime": kind.MIME.Value, "ext": kind.Extension}).WithError(err).Error("file type not supported")
			return errors.New("invalid file type")
		}
	}

	w, err := s.Bucket.NewWriter(ctx, fileName, nil)
	if err != nil {
		logrus.WithError(err).Error("unable to create upload writer")
		return errors.Wrap(err, "unable to creat upload writer")
	}

	defer w.Close()
	defer reader.Close()

	if _, err := io.Copy(w, reader); err != nil {
		logrus.WithError(err).Error("failed to upload file")
		return errors.Wrap(err, "failed to upload file")
	}

	return err
}
