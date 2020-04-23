package database

import (
	"alexandria/internal/backup"
	"alexandria/internal/documents"
	"alexandria/internal/journal"
	"alexandria/internal/links"
	"alexandria/internal/tags"
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func NewBackupRepository(psql *PostgresDatabase, neo *Neo4jDatabase) backup.Repository {
	return &backupRepo{
		postgres: psql,
		neo:      neo,
	}
}

type backupRepo struct {
	postgres *PostgresDatabase
	neo      *Neo4jDatabase
}

func (r *backupRepo) FindAllTags() ([]tags.Tag, error) {
	return r.postgres.FindAllTags()
}
func (r *backupRepo) FindAll(ctx context.Context, filter map[string]interface{}) ([]*documents.Document, error) {
	return r.postgres.FindAll(ctx, filter)
}
func (r *backupRepo) FindAllLinks() ([]links.Link, error) {
	return r.postgres.FindAllLinks()
}
func (r *backupRepo) FindAllEntries() ([]journal.Entry, error) {
	return r.postgres.FindAllEntries()
}
func (r *backupRepo) Restore(b backup.Backup) error {
	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		return r.postgres.Restore(b)
	})

	eg.Go(func() error {
		return r.neo.Restore(b)
	})

	if err := eg.Wait(); err != nil {
		logrus.WithError(err).Error("unable to restore")
		return err
	}
	return nil
}
