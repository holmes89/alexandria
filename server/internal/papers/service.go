package papers

import (
	"alexandria/internal/documents"
	"alexandria/internal/tags"
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"mime/multipart"
)

type PaperService interface {
	FindAll(ctx context.Context) ([]*documents.Document, error)
	FindByID(ctx context.Context, id string) (*documents.Document, error)
	Add(ctx context.Context, file multipart.File, paper *documents.Document) error
	AddTag(id string, tag string) error
	RemoveTag(id string, tag string) error
}

type service struct {
	docService documents.DocumentService
	tagsRepo   tags.Repository
}

func NewPaperService(docService documents.DocumentService, tagsRepo tags.Repository) PaperService {
	return &service{
		docService: docService,
		tagsRepo:   tagsRepo,
	}
}

func (s *service) FindAll(ctx context.Context) ([]*documents.Document, error) {
	m := map[string]interface{}{"type": "paper"}
	entities, err := s.docService.FindAll(ctx, m)
	if err != nil {
		logrus.WithError(err).Error("unable to fetch papers from repository")
		return nil, errors.Wrap(err, "unable to fetch from repository")
	}
	return entities, nil
}

func (s *service) FindByID(ctx context.Context, id string) (*documents.Document, error) {
	entity, err := s.docService.FindByID(ctx, id)
	if err != nil {
		logrus.WithError(err).WithField("id", id).Error("unable to fetch paper from repository")
		return nil, errors.Wrap(err, "unable to fetch from repository")
	}
	return entity, nil
}

func (s *service) Add(ctx context.Context, file multipart.File, paper *documents.Document) error {
	paper.Type = "paper"

	if err := s.docService.Add(ctx, file, paper); err != nil {
		logrus.WithError(err).Error("unable to save to repo")
		return errors.Wrap(err, "failed to store data in repo")
	}

	return nil
}

func (s *service) AddTag(id string, tag string) error {
	return s.tagsRepo.AddResourceTag(id, tags.PaperResource, tag)
}

func (s *service) RemoveTag(id string, tag string) error {
	return s.tagsRepo.RemoveResourceTag(id, tag)
}
