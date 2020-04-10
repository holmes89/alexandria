package backup

import (
	"alexandria/internal/common"
	"alexandria/internal/database"
	"alexandria/internal/documents"
	"alexandria/internal/journal"
	"alexandria/internal/links"
	"bytes"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
	"time"
)

type runner struct {
	ticker        *time.Ticker
	documentsRepo documents.DocumentRepository
	journalRepo   journal.Repository
	linksRepo     links.Repository
	storage       common.BackupSave
}

func NewBackupRunner(lc fx.Lifecycle, db *database.PostgresDatabase, storage common.BackupSave) {
	r := &runner{
		documentsRepo: db,
		journalRepo:   db,
		linksRepo:     db,
		storage:       storage,
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

type backup struct {
	Docs    []*documents.Document `json:"documents"`
	Journal []journal.Entry       `json:"journal_entries"`
	Links   []links.Link          `json:"links"`
}

func (r *runner) backup() {
	egroup, ctx := errgroup.WithContext(context.Background())

	b := &backup{}
	egroup.Go(func() error {
		docs, err := r.documentsRepo.FindAll(ctx, nil)
		b.Docs = docs
		return err
	})

	egroup.Go(func() error {
		entries, err := r.journalRepo.FindAllEntries()
		b.Journal = entries
		return err
	})

	egroup.Go(func() error {
		l, err := r.linksRepo.FindAllLinks()
		b.Links = l
		return err
	})

	if err := egroup.Wait(); err != nil {
		logrus.WithError(err).Fatal("unable to backup")
	}

	marshalled, err := json.Marshal(b)
	if err != nil {
		logrus.WithError(err).Fatal("unable to marshall backup")
	}

	location, err := r.storage.Save(context.Background(), "backup.json", bytes.NewReader(marshalled))
	if err != nil {
		logrus.WithError(err).Fatal("unable to send to bucket")
	}

	logrus.WithField("location", location).Info("backup successful")

}

func (r *runner) start() {
	go r.backup()                              // remove after testing
	r.ticker = time.NewTicker(8 * time.Minute) // 15 minute cool down so we want to backup at least once
	go func() {
		for {
			select {
			case <-r.ticker.C:
				go r.backup()
			}
		}
	}()
}

func (r *runner) stop() {
	r.ticker.Stop()
}
