package database

import "alexandria/internal/links"

type linksRepo struct {
	postgres *PostgresDatabase
	neo      *Neo4jDatabase
}

func NewLinksRepository(psql *PostgresDatabase, neo *Neo4jDatabase) links.Repository {
	return &linksRepo{
		postgres: psql,
		neo:      neo,
	}
}

func (r *linksRepo) FindAllLinks() ([]links.Link, error) {
	return r.FindAllLinks()
}

func (r *linksRepo) FindLinkByID(id string) (links.Link, error) {
	return r.postgres.FindLinkByID(id)
}

func (r *linksRepo) CreateLink(l links.Link) (links.Link, error) {
	nl, err := r.postgres.CreateLink(l)
	if err != nil {
		return l, err
	}
	return r.neo.CreateLink(nl)
}
