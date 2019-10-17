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
	return &PostgresDatabase{db}
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

func (r *PostgresDatabase) FindAll(ctx context.Context) (books []*Book, err error) {
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	rows, err := ps.Select("id", "description", "displayName", "name", "type", "path", "created", "modified").
		From("books").RunWith(r.conn).Query()

	for rows.Next() {
		book := &Book{}
		if err := rows.Scan(book.ID, book.Description, book.DisplayName, book.Name, book.Type, book.Path, book.Created, book.Modified); err != nil {
			logrus.WithError(err).Warn("unable to scan book results")
		}
		books = append(books, book)
	}
	return books, nil
}

func (r *PostgresDatabase) FindByID(ctx context.Context, id string) (*Book, error) {
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	row := ps.Select("id", "description", "displayName", "name", "type", "path", "created", "modified").
		From("books").Where(sq.Eq{"id": id}).RunWith(r.conn).QueryRow()

	book := &Book{}
	if err := row.Scan(book.ID, book.Description, book.DisplayName, book.Name, book.Type, book.Path, book.Created, book.Modified); err != nil {
		logrus.WithError(err).Warn("unable to scan book results")
	}

	return book, nil
}

func (r *PostgresDatabase) Insert(ctx context.Context, book *Book) error {
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	if _, err := ps.Insert("books").Columns("id", "description", "displayName", "name", "type", "path").
		Values(book.ID, book.Description, book.DisplayName, book.Name, book.Type, book.Path).
		RunWith(r.conn).
		Exec(); err != nil {
		logrus.WithError(err).Warn("unable to insert book")
		return errors.Wrap(err, "uanble to insert book metadata")
	}
	return nil
}
