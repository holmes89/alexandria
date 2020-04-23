package backup

import (
	"alexandria/internal/common"
	"alexandria/internal/documents"
	"alexandria/internal/journal"
	"alexandria/internal/links"
	"alexandria/internal/tags"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
	"time"
)

type service struct {
	backupRepo Repository
	storage    common.BackupStorage
}

type runner struct {
	service Service
	ticker  *time.Ticker
}

var fileNameBase = common.GetEnv("BACKUP_FILE", "backups/backup")

func NewBackupRunner(lc fx.Lifecycle, s Service) {
	r := &runner{
		service: s,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logrus.Info("starting server")
			go r.start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logrus.Info("stopping server")
			r.stop()
			return nil
		},
	})
}

type Service interface {
	Backup() error
	Restore(id string, restoreType Restore) error
}

type Repository interface {
	FindAllTags() ([]tags.Tag, error)
	FindAll(ctx context.Context, filter map[string]interface{}) ([]*documents.Document, error)
	FindAllLinks() ([]links.Link, error)
	FindAllEntries() ([]journal.Entry, error)
	Restore(b Backup) error
}

type Restore string

const (
	RestoreAll      Restore = "all"
	RestoreUnknown  Restore = "unknown"
	RestorePostgres Restore = "postgres"
	RestoreGraph    Restore = "graph"
)

func ParseRestore(req string) Restore {
	switch req {
	case "":
		return RestoreAll
	case "postgres":
		return RestorePostgres
	case "graph":
		return RestoreGraph
	default:
		return RestoreUnknown
	}
}

func NewService(db Repository, storage common.BackupStorage) Service {
	return &service{
		backupRepo: db,
		storage:    storage,
	}
}

type Backup struct {
	Docs    []*documents.Document `json:"documents"`
	Journal []journal.Entry       `json:"journal_entries"`
	Links   []links.Link          `json:"links"`
	Tags    []tags.Tag            `json:"tags"`
}

func (r *service) Restore(id string, restoreType Restore) error {

	fileName := fmt.Sprintf("%s-%s.json", fileNameBase, id)
	f, err := r.storage.Reader(context.Background(), fileName)
	if err != nil {
		logrus.WithError(err).Error("unable to fetch file")
		return errors.New("unable to download file")
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(f); err != nil {
		logrus.WithError(err).Error("unable to read file")
		return errors.New("unable to read file")
	}

	b := Backup{}
	if err := json.Unmarshal(buf.Bytes(), &b); err != nil {
		logrus.WithError(err).Error("unable to unmarshal file")
		return errors.New("unable to unmarshal file")
	}

	if err := r.backupRepo.Restore(b); err != nil {
		logrus.WithError(err).Error("unable to populate database")
		return errors.New("unable to populate database")
	}

	logrus.WithField("id", id).Info("backup restored")
	return nil
}
func (r *service) Backup() error {
	egroup, ctx := errgroup.WithContext(context.Background())

	b := &Backup{}
	egroup.Go(func() error {
		docs, err := r.backupRepo.FindAll(ctx, nil)
		b.Docs = docs
		return err
	})

	egroup.Go(func() error {
		entries, err := r.backupRepo.FindAllEntries()
		b.Journal = entries
		return err
	})

	egroup.Go(func() error {
		l, err := r.backupRepo.FindAllLinks()
		b.Links = l
		return err
	})

	egroup.Go(func() error {
		l, err := r.backupRepo.FindAllTags()
		b.Tags = l
		return err
	})

	if err := egroup.Wait(); err != nil {
		logrus.WithError(err).Error("unable to pull data from repositories")
		return errors.New("unable to pull data from repositories")
	}

	marshalled, err := json.Marshal(b)
	if err != nil {
		logrus.WithError(err).Error("unable to marshall Backup")
		return errors.New("unable to marshall Backup")
	}

	fileName := fmt.Sprintf("%s-%d.json", fileNameBase, time.Now().Unix())
	location, err := r.storage.Save(context.Background(), fileName, bytes.NewReader(marshalled))
	if err != nil {
		logrus.WithError(err).Error("unable to send to bucket")
		return errors.New("unable to persist file")
	}

	logrus.WithField("location", location).Info("Backup successful")
	return nil
}

func (r *runner) start() {
	r.ticker = time.NewTicker(1 * time.Hour)
	go func() {
		for {
			select {
			case <-r.ticker.C:
				go func() {
					if err := r.service.Backup(); err != nil {
						logrus.Fatal("unable to Backup")
					}
				}()
			}
		}
	}()
}

func (r *runner) stop() {
	r.ticker.Stop()
}
