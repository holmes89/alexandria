package database

import (
	"alexandria/internal/tags"
)

type tagsRepo struct {
	postgres *PostgresDatabase
	neo      *Neo4jDatabase
}

func NewTagsRepository(psql *PostgresDatabase, neo *Neo4jDatabase) tags.Repository {
	return &tagsRepo{
		postgres: psql,
		neo:      neo,
	}
}
func (r *tagsRepo) FindAllTags() ([]tags.Tag, error) {
	return r.postgres.FindAllTags()
}

func (r *tagsRepo) CreateTag(tag tags.Tag) (tags.Tag, error) {
	return r.postgres.CreateTag(tag)
}

func (r *tagsRepo) AddResourceTag(resourceID string, resourceType tags.ResourceType, tagName string) error {
	return r.postgres.AddResourceTag(resourceID, resourceType, tagName)
}

func (r *tagsRepo) RemoveResourceTag(resourceID string, tagName string) error {
	return r.postgres.RemoveResourceTag(resourceID, tagName)
}
