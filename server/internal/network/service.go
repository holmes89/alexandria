package network

import (
	"alexandria/internal/backup"
	"errors"
	"github.com/sirupsen/logrus"
)

type service struct {
	aggService backup.SystemAggregator
}

type Service interface {
	GetNetwork() (Network, error)
}

func NewService(aggService backup.SystemAggregator) Service {
	return &service{
		aggService: aggService,
	}
}

func (s *service) GetNetwork() (n Network, err error) {
	b, err := s.aggService.AggregateAllData()

	if err != nil {
		logrus.WithError(err).Error("unable to get aggregations")
		return n, errors.New("unable to aggregate data")
	}

	var nodes []Node
	var edges []Edge
	for _, d := range b.Tags {
		n, _ := createNodeAndEdges(d.ID, "tag", d.DisplayName, nil)
		nodes = append(nodes, n)
	}

	for _, d := range b.Docs {
		n, e := createNodeAndEdges(d.ID, d.Type, d.DisplayName, d.Tags)
		nodes = append(nodes, n)
		edges = append(edges, e...)
	}

	for _, d := range b.Links {
		n, e := createNodeAndEdges(d.ID, "link", d.DisplayName, d.Tags)
		nodes = append(nodes, n)
		edges = append(edges, e...)
	}

	n.Nodes = nodes
	n.Edges = edges

	return n, nil

}

func createNodeAndEdges(id, t, name string, tags []string) (Node, []Edge) {
	n := Node{
		ID:          id,
		DisplayName: name,
		Type:        t,
	}
	var edges []Edge
	for _, ta := range tags {
		edges = append(edges, Edge{
			NodeA: id,
			NodeB: ta,
		})
	}
	return n, edges
}
