package database

import (
	"alexandria/internal/common"
	"alexandria/internal/documents"
	"alexandria/internal/journal"
	"alexandria/internal/links"
	"alexandria/internal/user"
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

func NewPostgresDatabase(config common.PostgresDatabaseConfig) *PostgresDatabase {
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

func migrateDB(config common.PostgresDatabaseConfig) {
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

func (r *PostgresDatabase) FindAll(ctx context.Context, filter map[string]interface{}) (docs []*documents.Document, err error) {
	docs = []*documents.Document{}
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	rows, err := ps.Select("id", "description", "displayName", "name", "type", "path", "created", "updated").
		From("documents").Where(filter).RunWith(r.conn).Query()

	if err != nil {
		logrus.WithError(err).Error("unable to fetch results")
		return nil, errors.New("unable to fetch results")
	}
	for rows.Next() {
		doc := &documents.Document{}
		if err := rows.Scan(&doc.ID, &doc.Description, &doc.DisplayName, &doc.Name, &doc.Type, &doc.Path, &doc.Created, &doc.Updated); err != nil {
			logrus.WithError(err).Warn("unable to scan doc results")
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

func (r *PostgresDatabase) FindByID(ctx context.Context, id string) (*documents.Document, error) {
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	row := ps.Select("id", "description", "displayName", "name", "type", "path", "created", "updated").
		From("documents").Where(sq.Eq{"id": id}).RunWith(r.conn).QueryRow()
	doc := &documents.Document{}
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

func (r *PostgresDatabase) Insert(ctx context.Context, doc *documents.Document) error {
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

func (r *PostgresDatabase) UpsertStream(ctx context.Context, input <-chan *documents.Document) error {
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

func (r *PostgresDatabase) FindUserByUsername(ctx context.Context, username string) (*user.User, error) {
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	var entity user.User
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

func (r *PostgresDatabase) CreateUser(ctx context.Context, user *user.User) error {
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

func (r *PostgresDatabase) FindAllEntries() ([]journal.Entry, error) {
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	rows, err := ps.Select("id", "content", "created").From("journal_entry").RunWith(r.conn).Query()
	if err != nil {
		logrus.WithError(err).Error("unable to find entries")
		return nil, errors.New("unable to find entries")
	}
	entries := []journal.Entry{}
	for rows.Next() {
		var entry journal.Entry
		if err := rows.Scan(&entry.ID, &entry.Content, &entry.Created); err != nil {
			logrus.WithError(err).Warn("unable to scan entry")
		}
		entries = append(entries, entry)
	}
	rows.Close()
	return entries, nil
}

func (r *PostgresDatabase) CreateEntry(entry journal.Entry) (journal.Entry, error) {
	newEntry := journal.Entry{
		Content: entry.Content,
	}
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	if err := ps.Insert("journal_entry").
		Columns("content").
		Values(entry.Content).
		Suffix("RETURNING id, created").
		RunWith(r.conn).
		QueryRow().
		Scan(&newEntry.ID, &newEntry.Created); err != nil {

		logrus.WithError(err).Error("unable to insert entry")
		return newEntry, errors.New("unable to insert entry")
	}
	return newEntry, nil
}

func (r *PostgresDatabase) FindAllLinks() ([]links.Link, error) {
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	rows, err := ps.Select("id", "link", "display_name", "icon_path", "created").From("links").RunWith(r.conn).Query()
	if err != nil {
		logrus.WithError(err).Error("unable to find entries")
		return nil, errors.New("unable to find entries")
	}
	entries := []links.Link{}
	for rows.Next() {
		var entry links.Link
		if err := rows.Scan(&entry.ID, &entry.Link, &entry.DisplayName, &entry.IconPath, &entry.Created); err != nil {
			logrus.WithError(err).Warn("unable to scan entry")
		}
		entries = append(entries, entry)
	}
	rows.Close()
	return entries, nil
}

func (r *PostgresDatabase) CreateLink(entry links.Link) (links.Link, error) {
	newEntry := links.Link{
		Link:        entry.Link,
		DisplayName: entry.DisplayName,
		IconPath:    entry.IconPath,
	}
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	if err := ps.Insert("links").
		Columns("link", "display_name", "icon_path").
		Values(entry.Link, entry.DisplayName, entry.IconPath).
		Suffix("RETURNING id, created").
		RunWith(r.conn).
		QueryRow().
		Scan(&newEntry.ID, &newEntry.Created); err != nil {

		logrus.WithError(err).Error("unable to insert entry")
		return newEntry, errors.New("unable to insert entry")
	}
	return newEntry, nil
}

func NewPostgresDocumentRepository(database *PostgresDatabase) documents.DocumentRepository {
	return database
}

func NewUserPostgresRepository(database *PostgresDatabase) user.Repository {
	return database
}

func NewJournalRepository(database *PostgresDatabase) journal.Repository {
	return database
}

func NewLinksRepository(database *PostgresDatabase) links.Repository {
	return database
}
