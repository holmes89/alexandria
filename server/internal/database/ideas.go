package database

import (
	"alexandria/internal/ideas"
)

type ideasRepo struct {
	postgres *PostgresDatabase
	neo      *Neo4jDatabase
}

func NewIdeasRepository(psql *PostgresDatabase, neo *Neo4jDatabase) ideas.Repository {
	return &ideasRepo{
		postgres: psql,
		neo:      neo,
	}
}

func (r *ideasRepo) GetIdeas() ([]ideas.Idea, error) {

}

func (r *ideasRepo) GetIdeaByID(id string) (ideas.Idea, error) {

}

func (r *ideasRepo) CreateIdea(idea ideas.Idea) (ideas.Idea, error) {

}

func (r *ideasRepo) AddIdeaResource(resource ideas.IdeaResource) error {

}

func (r *ideasRepo) RemoveIdeaResource(resource ideas.IdeaResource) error {

}
