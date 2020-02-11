package main

import (
	"context"
	"crypto/tls"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)


var (
	ErrInvalidFileType = errors.New("invalid file type")
)

type Document struct {
	ID          string     `json:"id" firestore:"id"`
	DisplayName string     `json:"display_name" firestore:"display_name"`
	Name        string     `json:"name" firestore:"name"`
	Path        string     `json:"path" firestore:"path"`
	Type        string     `json:"type" firestore:"type"`
	Description string     `json:"description" firestore:"description"`
	Created     time.Time  `json:"created" firestore:"created"`
	Updated     *time.Time `json:"updated" firestore:"updated"`
}

type DocumentService interface {
	GetAll(ctx context.Context, filter map[string]interface{}) ([]*Document, error)
	GetByID(ctx context.Context, id string) (*Document, error)
	Add(ctx context.Context, file multipart.File, document *Document) error
	Delete(ctx context.Context, id string) error
	Scan(ctx context.Context) error
}

type DocumentRepository interface {
	FindAll(ctx context.Context, filter map[string]interface{}) ([]*Document, error)
	FindByID(ctx context.Context, id string) (*Document, error)
	Insert(ctx context.Context, document *Document) error
	Delete(ctx context.Context, id string) error
	UpsertStream(ctx context.Context, input <-chan *Document) error
}

func NewPostgresDocumentRepository(database *PostgresDatabase) DocumentRepository {
	return database
}

func NewFirestoreDocumentRepository(database *FirestoreDatabase) DocumentRepository {
	return database
}


type documentService struct {
	storage DocumentStorage
	repo    DocumentRepository
}

func NewDocumentService(storage DocumentStorage, repo DocumentRepository) DocumentService {
	return &documentService{
		storage: storage,
		repo:    repo,
	}
}

func (s *documentService) GetAll(ctx context.Context, filter map[string]interface{}) ([]*Document, error) {
	entities, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		logrus.WithError(err).Error("unable to fetch docs from repository")
		return nil, errors.Wrap(err, "unable to fetch from repository")
	}
	return entities, nil
}

func (s *documentService) GetByID(ctx context.Context, id string) (*Document, error) {
	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		logrus.WithError(err).WithField("id", id).Error("unable to fetch doc from repository")
		return nil, errors.Wrap(err, "unable to fetch from repository")
	}

	filePath, err := s.storage.Get(ctx, entity.Path)
	if err != nil {
		logrus.WithError(err).WithField("id", id).Error("unable to get path from storage")
		return nil, errors.Wrap(err, "unable to get path from storage")
	}
	entity.Path = filePath
	return entity, nil
}

func (s *documentService) Add(ctx context.Context, file multipart.File, doc *Document) error {
	if !isSupported(file) {
		return ErrInvalidFileType
	}
	path, err := s.storage.Save(ctx, doc.Name, file)
	if err != nil {
		logrus.WithError(err).Error("unable to write to storage")
		return errors.Wrap(err, "failed to write to storage")
	}

	doc.ID = uuid.New().String()
	doc.Path = path
	t := time.Now()
	doc.Created = t
	doc.Updated = &t

	if err := s.repo.Insert(ctx, doc); err != nil {
		logrus.WithError(err).Error("unable to save to repo")
		return errors.Wrap(err, "failed to store data in repo")
	}

	go s.CreateCover(doc.ID, doc.Path)
	return nil
}

func (s *documentService) CreateCover(id, path string) {
	url := getEnv("COVER_ENDPOINT", "")
	if url == "" {
		logrus.Panic("cover endpoint not set")
	}

	logrus.Infof("calling %s", url)
	if !strings.Contains(url, "http") {
		url = "https" + url
	}
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	resp, err := client.
		R().
		SetBody(coverRequest{ID:id, Path: path}).
		Post(url)

	if err != nil {
		logrus.WithError(err).Error("unable to create cover")
		return
	}
	if err := resp.Error(); err != nil {
		logrus.WithField("err", err).Error("unable to create cover")
		return
	}


	if resp.StatusCode() != http.StatusCreated {
		logrus.WithField("code", resp.StatusCode()).Error("request failed")
		return
	}

	logrus.Info("cover created")
}


type coverRequest struct {
	ID string `json:"id"`
	Path string `json:"path"`
}


func (s *documentService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *documentService) Scan(ctx context.Context) error {
	fileNameStream := s.storage.List(ctx)
	docStream := make(chan *Document)
	go func() {
		defer close(docStream)
		for path := range fileNameStream {
			ext := filepath.Ext(path)
			if ext != ".pdf" {
				continue
			}
			name := strings.ReplaceAll(path, ext, "")
			name = strings.ReplaceAll(name, filepath.Dir(path), "")
			if name[0] == '/' {
				name = name[1:]
			}
			doc := &Document{
				ID: uuid.New().String(),
				DisplayName: name,
				Name:        name,
				Path:        path,
				Type:        "book",
				Created:     time.Now(),
			}
			docStream <- doc
		}
	}()
	return s.repo.UpsertStream(ctx, docStream)
}

func isSupported(file multipart.File) bool {
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
