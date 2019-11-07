package main

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"mime/multipart"
)


type PaperService interface {
	GetAll(ctx context.Context) ([]*Document, error)
	GetByID(ctx context.Context, id string) (*Document, error)
	Add(ctx context.Context, file multipart.File, paper *Document) error
}

type paperService struct {
	docService DocumentService
}

func NewPaperService(service DocumentService) PaperService {
	return &paperService{
		docService: service,
	}
}

func (s *paperService) GetAll(ctx context.Context) ([]*Document, error) {
	m := map[string]interface{} { "type": "paper"}
	entities, err := s.docService.GetAll(ctx, m)
	if err != nil {
		logrus.WithError(err).Error("unable to fetch papers from repository")
		return nil, errors.Wrap(err, "unable to fetch from repository")
	}
	return entities, nil
}

func (s *paperService) GetByID(ctx context.Context, id string) (*Document, error) {
	entity, err := s.docService.GetByID(ctx, id)
	if err != nil {
		logrus.WithError(err).WithField("id", id).Error("unable to fetch paper from repository")
		return nil, errors.Wrap(err, "unable to fetch from repository")
	}
	return entity, nil
}

func (s *paperService) Add(ctx context.Context, file multipart.File, paper *Document) error {
	paper.Type = "paper"

	if err := s.docService.Add(ctx, file, paper); err != nil {
		logrus.WithError(err).Error("unable to save to repo")
		return errors.Wrap(err, "failed to store data in repo")
	}

	return nil
}