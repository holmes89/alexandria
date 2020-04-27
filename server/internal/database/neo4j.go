package database

import (
	"alexandria/internal/backup"
	"alexandria/internal/common"
	"alexandria/internal/documents"
	"alexandria/internal/links"
	"alexandria/internal/tags"
	"context"
	"errors"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"time"
)

type Neo4jDatabase struct {
	conn neo4j.Driver
}

func NewNeo4jDatabase(lc fx.Lifecycle, config common.Neo4jConfig) *Neo4jDatabase {
	logrus.Info("connecting to neo4j")

	var driver neo4j.Driver
	var err error
	if config.Username == "neo4j" { //Assume if this is the password its local
		driver, err = retryNeo4j(3, 10*time.Second, func() (driver neo4j.Driver, e error) {
			return neo4j.NewDriver(config.URI, neo4j.BasicAuth(config.Username, config.Password, ""), func(c *neo4j.Config) {
				c.Encrypted = false
			})
		})
	} else {
		driver, err = retryNeo4j(3, 10*time.Second, func() (driver neo4j.Driver, e error) {
			return neo4j.NewDriver(config.URI, neo4j.BasicAuth(config.Username, config.Password, ""))
		})
	}

	if err != nil {
		logrus.WithError(err).Fatal("unable to connect to neo4j")
	}
	logrus.Info("connected to neo4j")
	db := &Neo4jDatabase{driver}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logrus.Info("closing connection for neo4j")
			return driver.Close()
		},
	})
	return db
}
func retryNeo4j(attempts int, sleep time.Duration, callback func() (driver neo4j.Driver, e error)) (driver neo4j.Driver, e error) {
	for i := 0; i <= attempts; i++ {
		conn, err := callback()
		if err == nil {
			return conn, nil
		}
		time.Sleep(sleep)

		logrus.WithError(err).Error("error connecting to neo4j, retrying")
	}
	return nil, fmt.Errorf("after %d attempts, connection failed", attempts)
}

func (r *Neo4jDatabase) Restore(b backup.Backup) error {
	sess, err := r.conn.Session(neo4j.AccessModeWrite)
	if err != nil {
		logrus.WithError(err).Error("unable to create session")
		return errors.New("unable to create session")
	}
	defer sess.Close()

	// Create tag nodes
	for _, tag := range b.Tags {
		if _, err := r.CreateTag(tag); err != nil {
			logrus.WithError(err).Error("unable to create tag nodes")
			return errors.New("unable to create tag node")
		}
	}

	// Create Documents
	for _, doc := range b.Docs {
		resourceType := tags.BookResource
		if doc.Type == "paper" {
			resourceType = tags.PaperResource
		}
		if err := r.Insert(context.Background(), doc); err != nil {
			logrus.WithError(err).Error("unable to create document nodes")
			return errors.New("unable to create document node")
		}
		for _, tag := range doc.Tags {
			if err := r.addResourceTagByID(doc.ID, resourceType, tag); err != nil {
				logrus.WithError(err).Error("unable to create document tag edge")
				return errors.New("unable to create document tag edge")
			}
		}
	}

	// Links
	for _, link := range b.Links {
		if _, err := r.CreateLink(link); err != nil {
			logrus.WithError(err).Error("unable to create link nodes")
			return errors.New("unable to create link node")
		}
		for _, tag := range link.Tags {
			if err := r.addResourceTagByID(link.ID, tags.LinksResource, tag); err != nil {
				logrus.WithError(err).Error("unable to create link tag edge")
				return errors.New("unable to create link tag edge")
			}
		}
	}
	return nil
}

func (r *Neo4jDatabase) CreateLink(entity links.Link) (links.Link, error) {
	sess, err := r.conn.Session(neo4j.AccessModeWrite)
	if err != nil {
		logrus.WithError(err).Error("unable to create session")
		return entity, errors.New("unable to create session")
	}
	defer sess.Close()

	if _, err := sess.Run("CREATE (n:Link { id: $id, display_name: $display_name, link: $link, icon_path: $icon_path })", map[string]interface{}{
		"id":           entity.ID,
		"display_name": entity.DisplayName,
		"link":         entity.Link,
		"icon_path":    entity.IconPath,
	}); err != nil {
		logrus.WithError(err).Error("unable to create link nodes")
		return entity, errors.New("unable to create link node")
	}

	return entity, nil
}

func (r *Neo4jDatabase) Insert(ctx context.Context, entity *documents.Document) error {
	if entity.Type == "paper" {
		return r.insertPaper(ctx, entity)
	} else {
		return r.insertBook(ctx, entity)
	}

}

func (r *Neo4jDatabase) insertBook(_ context.Context, entity *documents.Document) error {
	sess, err := r.conn.Session(neo4j.AccessModeWrite)
	if err != nil {
		logrus.WithError(err).Error("unable to create session")
		return errors.New("unable to create session")
	}
	defer sess.Close()
	if _, err := sess.Run("CREATE (n:Book { id: $id, display_name: $display_name, path: $path, name: $name, description: $description })", map[string]interface{}{
		"id":           entity.ID,
		"display_name": entity.DisplayName,
		"path":         entity.Path,
		"name":         entity.Name,
		"description":  entity.Description,
	}); err != nil {
		logrus.WithError(err).Error("unable to create book nodes")
		return errors.New("unable to create book node")
	}
	return nil
}

func (r *Neo4jDatabase) insertPaper(_ context.Context, entity *documents.Document) error {
	sess, err := r.conn.Session(neo4j.AccessModeWrite)
	if err != nil {
		logrus.WithError(err).Error("unable to create session")
		return errors.New("unable to create session")
	}

	defer sess.Close()
	if _, err := sess.Run("CREATE (n:Paper { id: $id, display_name: $display_name, path: $path, name: $name, description: $description}) ", map[string]interface{}{
		"id":           entity.ID,
		"display_name": entity.DisplayName,
		"path":         entity.Path,
		"name":         entity.Name,
		"description":  entity.Description,
	}); err != nil {
		logrus.WithError(err).Error("unable to create paper nodes")
		return errors.New("unable to create paper node")
	}

	return nil
}

func (r *Neo4jDatabase) CreateTag(entity tags.Tag) (tags.Tag, error) {
	sess, err := r.conn.Session(neo4j.AccessModeWrite)
	if err != nil {
		logrus.WithError(err).Error("unable to create session")
		return entity, errors.New("unable to create session")
	}
	defer sess.Close()

	// Create tag nodes
	if _, err := sess.Run("MERGE (n:Tag { id: $id, display_name: $display_name, color: $color })", map[string]interface{}{
		"id":           entity.ID,
		"display_name": entity.DisplayName,
		"color":        entity.TagColor,
	}); err != nil {
		logrus.WithError(err).Error("unable to create tag nodes")
		return entity, errors.New("unable to create tag node")
	}

	return entity, nil
}

func (r *Neo4jDatabase) AddResourceTag(resourceID string, resourceType tags.ResourceType, tagName string) error {
	sess, err := r.conn.Session(neo4j.AccessModeWrite)
	if err != nil {
		logrus.WithError(err).Error("unable to create session")
		return errors.New("unable to create session")
	}
	defer sess.Close()

	nodeType := getNodeType(resourceType)
	cypher := fmt.Sprintf("MATCH (a:%s),(b:Tag) WHERE a.id = $resourceID AND b.display_name = $tagName CREATE (a)-[r:HAS_TAG]->(b)", nodeType)
	if _, err := sess.Run(cypher, map[string]interface{}{
		"resourceID": resourceID,
		"tagName":    tagName,
	}); err != nil {
		logrus.WithError(err).Error("unable to create tag edge")
		return errors.New("unable to create tag edge")
	}

	return nil
}

func (r *Neo4jDatabase) addResourceTagByID(resourceID string, resourceType tags.ResourceType, tagID string) error {
	sess, err := r.conn.Session(neo4j.AccessModeWrite)
	if err != nil {
		logrus.WithError(err).Error("unable to create session")
		return errors.New("unable to create session")
	}
	defer sess.Close()

	nodeType := getNodeType(resourceType)
	cypher := fmt.Sprintf("MATCH (a:%s),(b:Tag) WHERE a.id = $resourceID AND b.id = $tagID CREATE (a)-[r:HAS_TAG]->(b)", nodeType)
	if _, err := sess.Run(cypher, map[string]interface{}{
		"resourceID": resourceID,
		"tagID":      tagID,
	}); err != nil {
		logrus.WithError(err).Error("unable to create tag edge by id")
		return errors.New("unable to create tag edge by id")
	}

	return nil
}

func (r *Neo4jDatabase) RemoveResourceTag(resourceID string, tagName string) error {
	sess, err := r.conn.Session(neo4j.AccessModeWrite)
	if err != nil {
		logrus.WithError(err).Error("unable to create session")
		return errors.New("unable to create session")
	}
	defer sess.Close()

	if _, err := sess.Run("MATCH (a)-[r:HAS_TAG]->(b:Tag) WHERE a.id = $resourceID AND b.id = $tagID DELETE r", map[string]interface{}{
		"resourceID": resourceID,
		"tagName":    tagName,
	}); err != nil {
		logrus.WithError(err).Error("unable to delete tag edge")
		return errors.New("unable to delete tag edge")
	}

	return nil
}

func (r *Neo4jDatabase) GetTaggedResources(id string) ([]tags.TaggedResource, error) {
	tr := []tags.TaggedResource{}
	sess, err := r.conn.Session(neo4j.AccessModeWrite)
	if err != nil {
		logrus.WithError(err).Error("unable to create session")
		return tr, errors.New("unable to create session")
	}
	defer sess.Close()

	result, err := sess.Run("MATCH (a)-[r:HAS_TAG]->(b:Tag) WHERE b.id = $tagID RETURN a", map[string]interface{}{
		"tagID": id,
	})
	if err != nil {
		logrus.WithError(err).Error("unable to delete tag edge")
		return tr, errors.New("unable to delete tag edge")
	}

	for result.Next() {
		n, ok := result.Record().GetByIndex(0).(neo4j.Node)
		if !ok {
			logrus.Error("return not of node type")
			return tr, errors.New("not node type")
		}
		p := n.Props()
		t, err := getResourceType(n.Labels())
		if err != nil {
			logrus.WithError(err).Error("unable to find resource type")
			return tr, errors.New("unable to find resource type")
		}
		tr = append(tr, tags.TaggedResource{
			ResourceID:  p["id"].(string),
			DisplayName: p["display_name"].(string),
			Type:        t,
		})
	}
	return tr, nil
}

func getNodeType(resourceType tags.ResourceType) string {
	switch resourceType {
	case tags.LinksResource:
		return "Link"
	case tags.BookResource:
		return "Book"
	case tags.PaperResource:
		return "Paper"
	default:
		return "Unknown"
	}
}

func getResourceType(nodeTypes []string) (tags.ResourceType, error) {
	for _, nodeType := range nodeTypes {
		switch nodeType {
		case "Link":
			return tags.LinksResource, nil
		case "Book":
			return tags.BookResource, nil
		case "Paper":
			return tags.PaperResource, nil
		}
	}
	return "", errors.New("invalid resource type")
}
