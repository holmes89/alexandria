package main

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

type Service struct {
	bookBucket *BookBucket
	uploadBucket *SiteBucket
}

func NewService(bucket *BookBucket, siteBucket *SiteBucket) *Service {
	return &Service{
		bookBucket: bucket,
		uploadBucket: siteBucket,
	}
}

func (s *Service) CreateCover(id string, path string) error {
	//Download
	r, err := s.bookBucket.GetBook(path)
	if err != nil {
		logrus.WithError(err).Error("unable to fetch book")
		return errors.New("book doesn't exist")
	}

	file, err := ioutil.TempFile("/tmp", id)
	if _, err := io.Copy(file, r); err != nil {
		logrus.WithError(err).Error("unable to create temp file")
		return errors.New("unable to create temp file")
	}
	r.Close()

	logrus.WithField("id", id).Info("extracting cover")
	// Generate image
	//gs -sDEVICE=jpeg -dPDFFitPage=true -dDEVICEWIDTHPOINTS=250 -dDEVICEHEIGHTPOINTS=250 -sOutputFile=outputfile.jpeg inputfile.pdf
	args := []string{"-sDEVICE=jpeg","-dPDFFitPage=true","-dDEVICEWIDTHPOINTS=350","-dDEVICEHEIGHTPOINTS=350",fmt.Sprintf("-sOutputFile=/tmp/%s.jpg", id),file.Name()}

	cmd := exec.Command("gs", args...)
	if err := cmd.Run(); err != nil {
		logrus.WithError(err).Error("unable to create thumbnail")
		return errors.New("unable to create thumbnail")
	}

	logrus.WithField("id", id).Info("optimizing jpg")

	cmd = exec.Command("jpegoptim", fmt.Sprintf("/tmp/%s.jpg", id))
	if err := cmd.Run(); err != nil {
		logrus.WithError(err).Error("unable to optimize jpg")
		return errors.New("unable to optimize jpg")
	}

	if err := file.Close(); err != nil {
		logrus.WithError(err).Warn("unable to close file")
	}
	if err := os.Remove(file.Name()); err != nil {
		logrus.WithError(err).Warn("unable to delete file")
	}
	// Open file
	path = fmt.Sprintf("/tmp/%s.jpg", id)
	file, err = os.Open(path)
	if err != nil {
		logrus.WithError(err).Error("unable to open thumbnail")
		return errors.New("unable to open thumbnail")
	}

	// Upload
	if err := s.uploadBucket.UploadCover(id, file); err != nil {
		logrus.WithError(err).Error("unable to upload thumbnail")
		return errors.New("unable to upload thumbnail")
	}

	// Delete file
	if err := file.Close(); err != nil {
		logrus.WithError(err).Warn("unable to close file")
	}
	if err := os.Remove(file.Name()); err != nil {
		logrus.WithError(err).Warn("unable to delete file")
	}

	return nil
}
