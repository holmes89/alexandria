package database

import (
	"alexandria/internal/documents"
	"context"
)

type documentsRepo struct {
	postgres *PostgresDatabase
	neo      *Neo4jDatabase
}

func NewDocumentRepository(psql *PostgresDatabase, neo *Neo4jDatabase) documents.DocumentRepository {
	return &documentsRepo{
		postgres: psql,
		neo:      neo,
	}
}

func (r *documentsRepo) FindAll(ctx context.Context, filter map[string]interface{}) ([]*documents.Document, error) {
	return r.postgres.FindAll(ctx, filter)
}

func (r *documentsRepo) FindByID(ctx context.Context, id string) (*documents.Document, error) {
	return r.postgres.FindByID(ctx, id)
}

func (r *documentsRepo) Insert(ctx context.Context, document *documents.Document) error {
	if err := r.postgres.Insert(ctx, document); err != nil {
		return err
	}
	return r.neo.Insert(ctx, document)
}

func (r *documentsRepo) Delete(ctx context.Context, id string) error {
	return r.postgres.Delete(ctx, id)
}

func (r *documentsRepo) UpdateDocument(ctx context.Context, document documents.Document) (documents.Document, error) {
	return r.postgres.UpdateDocument(ctx, document)
}

func (r *documentsRepo) UpsertStream(ctx context.Context, input <-chan *documents.Document) error {
	return r.postgres.UpsertStream(ctx, input)
}
