package main

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"mime/multipart"
)

type BookService interface {
	GetAll(ctx context.Context) ([]*Document, error)
	GetByID(ctx context.Context, id string) (*Document, error)
	Add(ctx context.Context, file multipart.File, book *Document) error
}

type bookService struct {
	docService DocumentService
}

func NewBookService(service DocumentService) BookService {
	return &bookService{
		docService: service,
	}
}

func (s *bookService) GetAll(ctx context.Context) ([]*Document, error) {
	m := map[string]interface{}{"type": "book"}
	entities, err := s.docService.GetAll(ctx, m)
	if err != nil {
		logrus.WithError(err).Error("unable to fetch books from repository")
		return nil, errors.Wrap(err, "unable to fetch from repository")
	}
	return entities, nil
}

func (s *bookService) GetByID(ctx context.Context, id string) (*Document, error) {
	entity, err := s.docService.GetByID(ctx, id)
	if err != nil {
		logrus.WithError(err).WithField("id", id).Error("unable to fetch book from repository")
		return nil, errors.Wrap(err, "unable to fetch from repository")
	}
	return entity, nil
}

func (s *bookService) Add(ctx context.Context, file multipart.File, book *Document) error {
	book.Type = "book"

	if err := s.docService.Add(ctx, file, book); err != nil {
		logrus.WithError(err).Error("unable to save to repo")
		return errors.Wrap(err, "failed to store data in repo")
	}

	return nil
}
