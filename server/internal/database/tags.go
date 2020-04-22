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
	t, err := r.postgres.CreateTag(tag)
	if err != nil {
		return t, err
	}
	return r.neo.CreateTag(t)
}

func (r *tagsRepo) AddResourceTag(resourceID string, resourceType tags.ResourceType, tagName string) error {
	if err := r.postgres.AddResourceTag(resourceID, resourceType, tagName); err != nil {
		return err
	}
	// TODO create resource if not exists MERGE (p:Person{name:"Marina"})
	return r.neo.AddResourceTag(resourceID, resourceType, tagName)
}

func (r *tagsRepo) RemoveResourceTag(resourceID string, tagName string) error {
	if err := r.postgres.RemoveResourceTag(resourceID, tagName); err != nil {
		return err
	}
	return r.neo.RemoveResourceTag(resourceID, tagName)
}
