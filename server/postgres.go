package main

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq" // Used for specifying the type client we are creating
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
)

type PostgresDatabase struct {
	conn *sql.DB
}

func NewPostgresDatabase(config PostgresDatabaseConfig) *PostgresDatabase {
	logrus.Info("connecting to postgres")
	db, err := retryPostgres(3, 10*time.Second, func() (db *sql.DB, e error) {
		return sql.Open("postgres", config.ConnectionString)
	})
	if err != nil {
		logrus.WithError(err).Fatal("unable to connect to postgres")
	}
	logrus.Info("connected to postgres")
	psqldb := &PostgresDatabase{db}
	psqldb.initializeRepository()

	return psqldb
}

func (r *PostgresDatabase) initializeRepository() {
	query := `CREATE TABLE IF NOT EXISTS documents(
  				id VARCHAR(36) NOT NULL,
				description VARCHAR(1024),
  				displayName VARCHAR(255) NOT NULL,
  				name VARCHAR(255) NOT NULL,
  				type VARCHAR(255) NOT NULL,
				path VARCHAR(255) NOT NULL,
  				created timestamp NOT NULL DEFAULT current_timestamp,
  				updated timestamp NULL DEFAULT NULL,
  				PRIMARY KEY(id)
			);`
	if _, err := r.conn.Exec(query); err != nil {
		logrus.WithError(err).Fatal("unable to initialize database")
	}
}

func retryPostgres(attempts int, sleep time.Duration, callback func() (*sql.DB, error)) (*sql.DB, error) {
	for i := 0; i <= attempts; i++ {
		conn, err := callback()
		if err == nil {
			return conn, nil
		}
		time.Sleep(sleep)

		logrus.WithError(err).Error("error connecting to postgres, retrying")
	}
	return nil, fmt.Errorf("after %d attempts, connection failed", attempts)
}

func (r *PostgresDatabase) FindAll(ctx context.Context, filter map[string]interface{}) (docs []*Document, err error) {
	docs = []*Document{}
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	rows, err := ps.Select("id", "description", "displayName", "name", "type", "path", "created", "updated").
		From("documents").Where(filter).RunWith(r.conn).Query()

	if err != nil {
		logrus.WithError(err).Error("unable to fetch results")
		return nil, errors.Wrap(err, "unable to fetch results")
	}
	for rows.Next() {
		doc := &Document{}
		if err := rows.Scan(&doc.ID, &doc.Description, &doc.DisplayName, &doc.Name, &doc.Type, &doc.Path, &doc.Created, &doc.Updated); err != nil {
			logrus.WithError(err).Warn("unable to scan doc results")
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

func (r *PostgresDatabase) FindByID(ctx context.Context, id string) (*Document, error) {
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	row := ps.Select("id", "description", "displayName", "name", "type", "path", "created", "updated").
		From("documents").Where(sq.Eq{"id": id}).RunWith(r.conn).QueryRow()
	doc := &Document{}
	if err := row.Scan(&doc.ID, &doc.Description, &doc.DisplayName, &doc.Name, &doc.Type, &doc.Path, &doc.Created, &doc.Updated); err != nil {
		logrus.WithError(err).Warn("unable to scan doc results")
	}

	return doc, nil
}

func (r *PostgresDatabase) existsByPath(ctx context.Context, path string) (bool, error) {
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	row := ps.Select("count(id)").
		From("documents").Where(sq.Eq{"path": path}).RunWith(r.conn).QueryRow()
	var count int
	if err := row.Scan(&count); err != nil {
		logrus.WithError(err).Warn("unable to scan doc results")
	}

	return count > 0, nil
}

func (r *PostgresDatabase) Insert(ctx context.Context, doc *Document) error {
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	if _, err := ps.Insert("documents").Columns("id", "description", "displayName", "name", "type", "path").
		Values(doc.ID, doc.Description, doc.DisplayName, doc.Name, doc.Type, doc.Path).
		RunWith(r.conn).
		Exec(); err != nil {
		logrus.WithError(err).Warn("unable to insert doc")
		return errors.Wrap(err, "unable to insert doc metadata")
	}
	return nil
}

func (r *PostgresDatabase) UpsertStream(ctx context.Context, input <-chan *Document) error {
	count := 0
	for doc := range input {
		bctx := context.Background()
		if exists, _ := r.existsByPath(bctx, doc.Path); exists {
			continue
		}
		if err := r.Insert(bctx, doc); err != nil {
			logrus.WithError(err).Info("unable to upsert document")
			return errors.Wrap(err, "unable to upsert document")
		}
		count++
	}
	logrus.WithField("count", count).Info("documents added")
	return nil
}

func (r *PostgresDatabase) Delete(ctx context.Context, id string) error {
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	if _, err := ps.Delete("documents").Where(sq.Eq{"id": id}).RunWith(r.conn).Exec(); err != nil {
		logrus.WithError(err).Warn("unable to scan doc results")
		return errors.Wrap(err, "unable to delete")
	}

	return nil
}
