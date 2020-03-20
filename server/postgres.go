package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/lib/pq" // Used for specifying the type client we are creating
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
	migrateDB(config)

	return psqldb
}

func migrateDB(config PostgresDatabaseConfig) {
	db, err := sql.Open("postgres", config.ConnectionString)
	if err != nil {
		logrus.WithError(err).Fatal("unable to connect to postgres to migrate")
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logrus.WithError(err).Fatal("unable to get driver to migrate")
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mind", driver)
	if err != nil {
		logrus.WithError(err).Fatal("unable to create migration instance")
	}
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			logrus.WithError(err).Fatal("unable to migrate")
		}
		logrus.Info("no migrations to run")
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
		return nil, errors.New( "unable to fetch results")
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
		return errors.New("unable to insert doc metadata")
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
			return errors.New("unable to upsert document")
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
		return errors.New("unable to delete")
	}

	return nil
}

func (r *PostgresDatabase) FindUserByUsername(ctx context.Context, username string) (*User, error) {
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	var entity User
	if err := ps.Select("id", "username", "password").
		From("users").
		Where(sq.Eq{"username": username}).
		RunWith(r.conn).
		QueryRow().
		Scan(&entity.ID, &entity.Username, &entity.Password); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			logrus.WithError(err).Error("could not find user")
			return nil, errors.New("could not find user")
	}
	return &entity, nil
}

func (r *PostgresDatabase) CreateUser(ctx context.Context, user *User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	if _, err := ps.Insert("users").
		Columns("id", "username", "password").
		Values(user.ID, user.Username, user.Password).
		RunWith(r.conn).Exec(); err != nil {

		logrus.WithError(err).Error("unable to create user")
		return errors.New("unable to create user")
	}
	return nil
}
