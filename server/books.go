package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"time"
)

var (
	ErrInvalidFileType = errors.New("invalid file type")
)

type Document struct {
	ID          string     `json:"id"`
	DisplayName string     `json:"display_name"`
	Name        string     `json:"name"`
	Path        string     `json:"path"`
	Type        string     `json:"type"`
	Description string     `json:"description"`
	Created     time.Time  `json:"created"`
	Updated     *time.Time `json:"updated"`
}

type BookService interface {
	GetAll(ctx context.Context) ([]*Document, error)
	GetByID(ctx context.Context, id string) (*Document, error)
	Add(ctx context.Context, file multipart.File, book *Document) error
}

type BookRepository interface {
	FindAll(ctx context.Context) ([]*Document, error)
	FindByID(ctx context.Context, id string) (*Document, error)
	Insert(ctx context.Context, book *Document) error
}

func NewPostgresBookRepository(database *PostgresDatabase) BookRepository {
	return database
}

type bookService struct {
	storage BookSave
	repo    BookRepository
}

func NewBookService(storage BookSave, repo BookRepository) BookService {
	return &bookService{
		storage: storage,
		repo:    repo,
	}
}

func (s *bookService) GetAll(ctx context.Context) ([]*Document, error) {
	entities, err := s.repo.FindAll(ctx)
	if err != nil {
		logrus.WithError(err).Error("unable to fetch books from repository")
		return nil, errors.Wrap(err, "unable to fetch from repository")
	}
	return entities, nil
}

func (s *bookService) GetByID(ctx context.Context, id string) (*Document, error) {
	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		logrus.WithError(err).WithField("id", id).Error("unable to fetch book from repository")
		return nil, errors.Wrap(err, "unable to fetch from repository")
	}
	return entity, nil
}

func (s *bookService) Add(ctx context.Context, file multipart.File, book *Document) error {
	if !isBook(file) {
		return ErrInvalidFileType
	}
	path, err := s.storage.Save(ctx, book.Name, file)
	if err != nil {
		logrus.WithError(err).Error("unable to write to storage")
		return errors.Wrap(err, "failed to write to storage")
	}

	book.ID = uuid.New().String()
	book.Path = path
	t := time.Now()
	book.Created = t
	book.Updated = &t
	book.Type = "book"

	if err := s.repo.Insert(ctx, book); err != nil {
		logrus.WithError(err).Error("unable to save to repo")
		return errors.Wrap(err, "failed to store data in repo")
	}

	return nil
}

func isBook(file multipart.File) bool {
	head := make([]byte, 261)
	if bytesRead, err := io.ReadFull(file, head); err == io.EOF {
		logrus.WithField("bytesRead", bytesRead).WithError(err).Error("couldn't read file header: unexpected EOF")
		return false
	} else if err != nil {
		logrus.WithField("bytesRead", bytesRead).WithError(err).Error("couldn't read file header")
		return false
	}

	file.Seek(0, io.SeekStart)

	if kind, err := filetype.Match(head); err != nil {
		logrus.WithError(err).Error("unable to determine file type")
		return false
	} else {
		if kind != matchers.TypeEpub && kind != matchers.TypePdf {
			logrus.WithFields(logrus.Fields{"mime": kind.MIME.Value, "ext": kind.Extension}).WithError(err).Error("file type not supported")
			return false
		}
	} // TODO mobi check
	return true
}