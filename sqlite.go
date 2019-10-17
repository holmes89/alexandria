package main

import (
	"context"
	"crawshaw.io/sqlite"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type SQLiteDatabase struct {
	pool *sqlite.Pool
}

func NewSQLiteDatabase(config SQLiteDatabaseConfig) *SQLiteDatabase {
	filename := fmt.Sprintf("file:%s", config.File)
	dbpool, err := sqlite.Open(filename, 0, 10)
	if err != nil {
		logrus.WithError(err).Fatal("unable to create db connection pool")
		return nil
	}
	return &SQLiteDatabase{
		pool: dbpool,
	}
}

func (r *SQLiteDatabase) FindAll(ctx context.Context) (books []*Book, err error) {
	conn := r.pool.Get(ctx.Done())
	if conn == nil {
		logrus.Error("no connections available in pool")
		return nil, errors.New("no connection available")
	}

	defer r.pool.Put(conn)

	stmt := conn.Prep("SELECT * FROM books;")
	for {
		if hasRow, err := stmt.Step(); err != nil {
			logrus.WithError(err).Error("statement creation failed")
			return nil, err
		} else if !hasRow {
			break
		}

		book := r.bookFromStatement(stmt)
		books = append(books, book)
	}

	return books, nil
}

func (r *SQLiteDatabase) FindByID(ctx context.Context, id string) (*Book, error) {
	conn := r.pool.Get(ctx.Done())
	if conn == nil {
		logrus.Error("no connections available in pool")
		return nil, errors.New("no connection available")
	}

	defer r.pool.Put(conn)

	stmt := conn.Prep("SELECT * FROM books where id = $id;")
	stmt.SetText("$id", id)

	var book *Book
	for {
		if hasRow, err := stmt.Step(); err != nil {
			logrus.WithField("err", err.Error()).Error("statement creation failed")
			return nil, err
		} else if !hasRow {
			break
		}
		book = r.bookFromStatement(stmt)
	}
	return book, nil
}

func (r *SQLiteDatabase) Insert(ctx context.Context, book *Book) error {
	conn := r.pool.Get(ctx.Done())
	if conn == nil {
		logrus.Error("no connections available in pool")
		return errors.New("no connection available")
	}

	defer r.pool.Put(conn)

	stmt := conn.Prep("INSERT INTO books VALUES($id, $displayName, $name, $path, $type, $description, $created, $modified);")
	stmt.SetText("$id", book.ID)
	stmt.SetText("$displayName", book.DisplayName)
	stmt.SetText("$name", book.Name)
	stmt.SetText("$path", book.Path)
	stmt.SetText("$type", book.Type)
	stmt.SetText("$description", book.Description)
	stmt.SetText("$created", book.Created.Format(time.RFC3339))
	stmt.SetText("$modified", book.Modified.Format(time.RFC3339))

	if _, err := stmt.Step(); err != nil {
		logrus.WithError(err).Error("unable to create book")
		return errors.New("not able to insert book")
	}

	return nil
}

func (r *SQLiteDatabase) bookFromStatement(stmt *sqlite.Stmt) *Book {
	created, err := time.Parse(time.RFC3339, stmt.GetText("created"))
	if err != nil {
		logrus.WithError(err).Warn("unable to parse date")
		created = time.Time{}
	}
	mod, err := time.Parse(time.RFC3339, stmt.GetText("modified"))
	if err != nil {
		logrus.WithError(err).Warn("unable to parse date")
		mod = time.Time{}
	}
	return &Book{
		ID:          stmt.GetText("id"),
		DisplayName: stmt.GetText("display_name"),
		Name:        stmt.GetText("name"),
		Path:        stmt.GetText("path"),
		Type:        stmt.GetText("type"),
		Description: stmt.GetText("description"),
		Created:     created,
		Modified:    mod,
	}
}
