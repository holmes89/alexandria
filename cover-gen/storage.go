package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gocloud.dev/blob"
	"gocloud.dev/blob/gcsblob"
	"gocloud.dev/gcp"
	"io"
	"net/url"
)

type BookBucket struct {
	bucket *blob.Bucket
}

func NewBookBucket() *BookBucket {
	config := LoadBookBucketConfig()
	return &BookBucket{
		bucket: NewGCPBucketStorage(config.BucketConfig),
	}
}

func (r *BookBucket) GetBook(name string) (io.ReadCloser, error) {
	logrus.WithField("name", name).Info("downloading")
	reader, err := r.bucket.NewReader(context.Background(), name, nil)
	if err != nil {
		logrus.WithError(err).Error("unable to fetch book")
		return nil, errors.New("unable to fetch book")
	}
	return reader, nil
}

type SiteBucket struct {
	bucket *blob.Bucket
}


func NewSiteBucket() *SiteBucket {
	config := LoadSiteBucketConfig()
	return &SiteBucket{
		bucket: NewGCPBucketStorage(config.BucketConfig),
	}
}

func (r *SiteBucket) UploadCover(id string, file io.Reader) error {
	logrus.WithField("id", id).Info("uploading")
	path := fmt.Sprintf("assets/covers/%s.jpg", id)
	writer, err := r.bucket.NewWriter(context.Background(), path, nil)
	if err != nil {
		logrus.WithError(err).Error("unable to create writer")
		return errors.New("unable to create writer")
	}
	defer writer.Close()
	if _, err := io.Copy(writer, file); err != nil {
		logrus.WithError(err).Error("unable to upload file")
		return errors.New("unable to upload file")
	}
	return nil
}



func NewGCPBucketStorage(config BucketConfig) *blob.Bucket {
	ctx := context.Background()

	urlString := config.ConnectionString
	urlParts, _ := url.Parse(urlString)
	// Your GCP credentials.
	// See https://cloud.google.com/docs/authentication/production
	// for more info on alternatives.
	creds, err := gcp.DefaultCredentials(ctx)
	if err != nil {
		logrus.Fatal(err)
	}

	accessID := config.AccessID
	accessKey := config.AccessKey

	if accessID == "" || accessKey == "" {
		logrus.Warn("unable to find access information using default credentials")
		credsMap := make(map[string]string)
		json.Unmarshal(creds.JSON, &credsMap)
		accessID = credsMap["client_id"]
		accessKey = credsMap["private_key"]
	}

	opts := &gcsblob.Options{
		GoogleAccessID: accessID,
		PrivateKey: []byte(accessKey),
	}
	// Create an HTTP client.
	// This example uses the default HTTP transport and the credentials
	// created above.
	client, err := gcp.NewHTTPClient(
		gcp.DefaultTransport(),
		gcp.CredentialsTokenSource(creds))
	if err != nil {
		logrus.Fatal(err)
	}

	// Create a *blob.Bucket.
	bucket, err := gcsblob.OpenBucket(ctx, client, urlParts.Host, opts)
	if err != nil {
		logrus.Fatal(err)
	}
	return bucket
}
