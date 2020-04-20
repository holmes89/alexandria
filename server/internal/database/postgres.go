package database

import (
	"alexandria/internal/backup"
	"alexandria/internal/common"
	"alexandria/internal/documents"
	"alexandria/internal/journal"
	"alexandria/internal/links"
	"alexandria/internal/tags"
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
	"github.com/iancoleman/strcase"
	_ "github.com/lib/pq" // Used for specifying the type client we are creating
	"github.com/sirupsen/logrus"
	"strings"
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
	rows, err := ps.Select("documents.id", "description", "display_name", "name", "type", "path", "COALESCE(string_agg(tagged_resources.id::character varying, ','), '')", "created", "updated").
		From("documents").
		LeftJoin("tagged_resources ON documents.id=tagged_resources.resource_id").
		Suffix("GROUP BY documents.id ORDER BY display_name ASC").
		Where(filter).RunWith(r.conn).Query()

	if err != nil {
		logrus.WithError(err).Error("unable to fetch results")
		return nil, errors.New("unable to fetch results")
	}
	for rows.Next() {
		doc := &documents.Document{}
		var tagList string
		doc.Tags = []string{}
		if err := rows.Scan(&doc.ID, &doc.Description, &doc.DisplayName, &doc.Name, &doc.Type, &doc.Path, &tagList, &doc.Created, &doc.Updated); err != nil {
			logrus.WithError(err).Warn("unable to scan doc results")
		}
		if tagList != "" {
			doc.Tags = append(doc.Tags, strings.Split(tagList, ",")...)
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

func (r *PostgresDatabase) FindByID(ctx context.Context, id string) (*documents.Document, error) {
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	row := ps.Select("documents.id", "description", "display_name", "name", "type", "path", "COALESCE(string_agg(tagged_resources.id::character varying, ','), '')", "created", "updated").
		From("documents").
		LeftJoin("tagged_resources ON documents.id=tagged_resources.resource_id").
		Suffix("GROUP BY documents.id").
		Where(sq.Eq{"documents.id": id}).RunWith(r.conn).QueryRow()
	doc := &documents.Document{}
	var tagList string
	doc.Tags = []string{}
	if err := row.Scan(&doc.ID, &doc.Description, &doc.DisplayName, &doc.Name, &doc.Type, &doc.Path, &tagList, &doc.Created, &doc.Updated); err != nil {
		logrus.WithError(err).Warn("unable to scan doc results")
	}
	if tagList != "" {
		doc.Tags = append(doc.Tags, strings.Split(tagList, ",")...)
	}

	return doc, nil
}

func (r *PostgresDatabase) UpdateDocument(_ context.Context, doc documents.Document) (result documents.Document, err error) {
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	_, err = ps.Update("documents").SetMap(
		map[string]interface{}{
			"description":  doc.Description,
			"display_name": doc.DisplayName,
			"type":         doc.Type,
			"updated":      time.Now()}).
		Where(sq.Eq{"id": doc.ID}).RunWith(r.conn).Exec()

	if err != nil {
		logrus.WithError(err).Error("unable to update doc")
		return result, errors.New("unable to update doc")
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
	if _, err := ps.Insert("documents").Columns("id", "description", "display_name", "name", "type", "path").
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
	rows, err := ps.Select("links.id", "link", "display_name", "icon_path", "COALESCE(string_agg(tagged_resources.id::character varying, ','), '')", "created").
		From("links").
		LeftJoin("tagged_resources ON links.id=tagged_resources.resource_id").Suffix("GROUP BY links.id ORDER BY created DESC").RunWith(r.conn).Query()
	if err != nil {
		logrus.WithError(err).Error("unable to find links")
		return nil, errors.New("unable to find links")
	}
	entries := []links.Link{}
	for rows.Next() {
		var entry links.Link
		var tagList string
		entry.Tags = []string{}
		if err := rows.Scan(&entry.ID, &entry.Link, &entry.DisplayName, &entry.IconPath, &tagList, &entry.Created); err != nil {
			logrus.WithError(err).Warn("unable to scan link")
		}
		if tagList != "" {
			entryTags := strings.Split(tagList, ",")
			entry.Tags = append(entry.Tags, entryTags...)
		}
		entries = append(entries, entry)
	}
	rows.Close()
	return entries, nil
}

func (r *PostgresDatabase) FindLinkByID(id string) (entity links.Link, err error) {
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	rowscanner := ps.Select("links.id", "link", "display_name", "icon_path", "COALESCE(string_agg(tagged_resources.id::character varying, ','), '')", "created").
		From("links").
		LeftJoin("tagged_resources ON links.id=tagged_resources.resource_id").
		Where(sq.Eq{"links.id": id}).
		Suffix("GROUP BY links.id ORDER BY created DESC").RunWith(r.conn).QueryRow()
	if err != nil {
		logrus.WithError(err).Error("unable to find links")
		return entity, errors.New("unable to find links")
	}
	var entry links.Link
	var tagList string
	entry.Tags = []string{}
	if err := rowscanner.Scan(&entry.ID, &entry.Link, &entry.DisplayName, &entry.IconPath, &tagList, &entry.Created); err != nil {
		logrus.WithError(err).Warn("unable to scan link")
	}
	if tagList != "" {
		entryTags := strings.Split(tagList, ",")
		entry.Tags = append(entry.Tags, entryTags...)
	}
	return entry, nil
}

func (r *PostgresDatabase) CreateLink(entry links.Link) (links.Link, error) {
	newEntry := links.Link{
		Link:        entry.Link,
		DisplayName: entry.DisplayName,
		IconPath:    entry.IconPath,
		Tags:        []string{},
	}
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	if err := ps.Insert("links").
		Columns("link", "display_name", "icon_path").
		Values(entry.Link, entry.DisplayName, entry.IconPath).
		Suffix("RETURNING id, created").
		RunWith(r.conn).
		QueryRow().
		Scan(&newEntry.ID, &newEntry.Created); err != nil {

		logrus.WithError(err).Error("unable to insert link")
		return newEntry, errors.New("unable to insert link")
	}
	return newEntry, nil
}

func (r *PostgresDatabase) FindAllTags() ([]tags.Tag, error) {
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	rows, err := ps.Select("id", "display_name", "color").From("tags").Suffix("ORDER BY display_name ASC").RunWith(r.conn).Query()
	if err != nil {
		logrus.WithError(err).Error("unable to find tags")
		return nil, errors.New("unable to find tags")
	}
	entries := []tags.Tag{}
	for rows.Next() {
		var entry tags.Tag
		if err := rows.Scan(&entry.ID, &entry.DisplayName, &entry.TagColor); err != nil {
			logrus.WithError(err).Warn("unable to scan tags")
		}
		entries = append(entries, entry)
	}
	rows.Close()
	return entries, nil
}

func (r *PostgresDatabase) CreateTag(entry tags.Tag) (tags.Tag, error) {
	color := tags.GetRandomColor()
	newEntry := tags.Tag{
		DisplayName: strcase.ToKebab(entry.DisplayName),
		TagColor:    color,
	}
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	if err := ps.Insert("tags").
		Columns("display_name", "color").
		Values(newEntry.DisplayName, color).
		Suffix("ON CONFLICT DO NOTHING RETURNING id").
		RunWith(r.conn).
		QueryRow().
		Scan(&newEntry.ID); err != nil {
		if err == sql.ErrNoRows {
			return newEntry, nil
		}
		logrus.WithError(err).Error("unable to insert tag")
		return newEntry, errors.New("unable to insert tag")
	}
	return newEntry, nil
}

func (r *PostgresDatabase) AddResourceTag(resourceID string, resourceType tags.ResourceType, tagName string) error {
	tagName = strcase.ToKebab(tagName)
	if _, err := r.CreateTag(tags.Tag{DisplayName: tagName}); err != nil {
		logrus.WithError(err).Error("unable to upsert tag")
		return errors.New("unable to upsert tag")
	}
	if _, err := r.conn.Exec("INSERT INTO tagged_resources(id, resource_id, resource_type) SELECT id, $1, $2  FROM tags where display_name = $3", resourceID, resourceType, tagName); err != nil {
		logrus.WithError(err).Error("unable to add tag")
		return errors.New("unable to add tag")
	}
	return nil
}

func (r *PostgresDatabase) RemoveResourceTag(resourceID string, tagName string) error {
	tagName = strcase.ToKebab(tagName)
	if _, err := r.conn.Exec("DELETE FROM tagged_resources where resource_id = $1 and id = (SELECT id FROM tags where display_name = $2)", resourceID, tagName); err != nil {
		logrus.WithError(err).Error("unable to remove tag")
		return errors.New("unable to remove tag")
	}
	return nil
}

func (r *PostgresDatabase) Restore(b backup.Backup) error {

	tx, err := r.conn.BeginTx(context.Background(), nil)
	if err != nil {
		logrus.WithError(err).Error("unable to create transaction")
		return errors.New("unable to create transaction")
	}

	if err := r.bulkInsertTags(tx, b.Tags); err != nil {
		logrus.WithError(err).Error("unable to insert tags")
		tx.Rollback()
		return errors.New("unable to insert tags")
	}
	if err := r.bulkInsertEntries(tx, b.Journal); err != nil {
		logrus.WithError(err).Error("unable to insert entries")
		tx.Rollback()
		return errors.New("unable to insert entries")
	}

	docsTrs, err := r.bulkInsertDocuments(tx, b.Docs)
	if err != nil {
		logrus.WithError(err).Error("unable to insert tags")
		tx.Rollback()
		return errors.New("unable to insert tags")
	}

	linksTrs, err := r.bulkInsertLinks(tx, b.Links)
	if err != nil {
		logrus.WithError(err).Error("unable to insert tags")
		tx.Rollback()
		return errors.New("unable to insert tags")
	}

	trs := append(docsTrs, linksTrs...)
	if err := r.bulkInsertTaggedResources(tx, trs); err != nil {
		logrus.WithError(err).Error("unable to insert tags")
		tx.Rollback()
		return errors.New("unable to insert tags")
	}

	if err := tx.Commit(); err != nil {
		logrus.WithError(err).Error("unable to commit transaction")
		return errors.New("unable to commit transaction")
	}
	return nil
}

func (r *PostgresDatabase) bulkInsertTags(tx *sql.Tx, tags []tags.Tag) error {
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	s := ps.Insert("tags").Columns("id", "display_name", "color")

	for _, t := range tags {
		s = s.Values(t.ID, t.DisplayName, t.TagColor)
	}

	if _, err := s.RunWith(tx).Exec(); err != nil {
		logrus.WithError(err).Error("unable to insert tags")
		return errors.New("unable to insert tags")
	}

	return nil
}

func (r *PostgresDatabase) bulkInsertTaggedResources(tx *sql.Tx, tags []tags.TaggedResource) error {
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	s := ps.Insert("tagged_resources").Columns("id", "resource_id", "resource_type")

	for _, t := range tags {
		s = s.Values(t.ID, t.ResourceID, t.Type)
	}

	if _, err := s.RunWith(tx).Exec(); err != nil {
		logrus.WithError(err).Error("unable to insert tagged resources")
		return errors.New("unable to insert tagged resources")
	}

	return nil
}

func (r *PostgresDatabase) bulkInsertDocuments(tx *sql.Tx, docs []*documents.Document) (tr []tags.TaggedResource, err error) {
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	s := ps.Insert("documents").Columns("id", "description", "display_name", "name", "type", "path")
	for _, d := range docs {

		for _, t := range d.Tags {
			ty := tags.BookResource
			if d.Type == "paper" {
				ty = tags.PaperResource
			}
			tr = append(tr, tags.TaggedResource{
				ID:         t,
				ResourceID: d.ID,
				Type:       ty,
			})
		}

		s = s.Values(d.ID, d.Description, d.DisplayName, d.Name, d.Type, d.Path)
	}

	if _, err := s.RunWith(tx).Exec(); err != nil {
		logrus.WithError(err).Error("unable to insert docs")
		return tr, errors.New("unable to insert docs")
	}

	return tr, nil
}

func (r *PostgresDatabase) bulkInsertEntries(tx *sql.Tx, entries []journal.Entry) error {
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	s := ps.Insert("journal_entry").Columns("id", "content")

	for _, e := range entries {
		s = s.Values(e.ID, e.Content)
	}

	if _, err := s.RunWith(tx).Exec(); err != nil {
		logrus.WithError(err).Error("unable to insert entries")
		return errors.New("unable to insert entries")
	}

	return nil
}

func (r *PostgresDatabase) bulkInsertLinks(tx *sql.Tx, lks []links.Link) (tr []tags.TaggedResource, err error) {
	ps := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	s := ps.Insert("links").Columns("id", "link", "icon_path", "display_name")
	for _, l := range lks {

		for _, t := range l.Tags {
			tr = append(tr, tags.TaggedResource{
				ID:         t,
				ResourceID: l.ID,
				Type:       tags.LinksResource,
			})
		}

		s = s.Values(l.ID, l.Link, l.IconPath, l.DisplayName)
	}

	if _, err := s.RunWith(tx).Exec(); err != nil {
		logrus.WithError(err).Error("unable to insert links")
		return tr, errors.New("unable to insert links")
	}

	return tr, nil
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

func NewTagsRepository(database *PostgresDatabase) tags.Repository {
	return database
}

func NewBackupRepository(database *PostgresDatabase) backup.Repository {
	return database
}
