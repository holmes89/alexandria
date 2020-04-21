package database

import (
	"alexandria/internal/backup"
	"alexandria/internal/common"
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

	driver, err := retryNeo4j(3, 10*time.Second, func() (driver neo4j.Driver, e error) {
		return neo4j.NewDriver(config.URI, neo4j.BasicAuth(config.Username, config.Password, ""), func(c *neo4j.Config) {
			c.Encrypted = false
		})
	})
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
		if _, err := sess.Run("CREATE (n:Tag { id: $id, display_name: $display_name, color: $color })", map[string]interface{}{
			"id":           tag.ID,
			"display_name": tag.DisplayName,
			"color":        tag.TagColor,
		}); err != nil {
			logrus.WithError(err).Error("unable to create tag nodes")
			return errors.New("unable to create tag node")
		}
	}

	// Create Documents
	for _, doc := range b.Docs {
		if doc.Type == "paper" {
			if _, err := sess.Run("CREATE (n:Paper { id: $id, display_name: $display_name, path: $path, name: $name, description: $description}) ", map[string]interface{}{
				"id":           doc.ID,
				"display_name": doc.DisplayName,
				"path":         doc.Path,
				"name":         doc.Name,
				"description":  doc.Description,
			}); err != nil {
				logrus.WithError(err).Error("unable to create paper nodes")
				return errors.New("unable to create paper node")
			}
			for _, tag := range doc.Tags {
				if _, err := sess.Run("MATCH (a:Paper),(b:Tag) WHERE a.id = $paperID AND b.id = $tagID CREATE (a)-[r:HAS_TAG]->(b)", map[string]interface{}{
					"paperID": doc.ID,
					"tagID":   tag,
				}); err != nil {
					logrus.WithError(err).Error("unable to create paper tag edge")
					return errors.New("unable to create paper tag edge")
				}
			}
		} else {
			if _, err := sess.Run("CREATE (n:Book { id: $id, display_name: $display_name, path: $path, name: $name, description: $description })", map[string]interface{}{
				"id":           doc.ID,
				"display_name": doc.DisplayName,
				"path":         doc.Path,
				"name":         doc.Name,
				"description":  doc.Description,
			}); err != nil {
				logrus.WithError(err).Error("unable to create book nodes")
				return errors.New("unable to create book node")
			}
			for _, tag := range doc.Tags {
				if _, err := sess.Run("MATCH (a:Book),(b:Tag) WHERE a.id = $paperID AND b.id = $tagID CREATE (a)-[r:HAS_TAG]->(b)", map[string]interface{}{
					"paperID": doc.ID,
					"tagID":   tag,
				}); err != nil {
					logrus.WithError(err).Error("unable to create book tag edge")
					return errors.New("unable to create book tag edge")
				}
			}
		}
	}

	// Links
	for _, link := range b.Links {
		if _, err := sess.Run("CREATE (n:Link { id: $id, display_name: $display_name, link: $link, icon_path: $icon_path })", map[string]interface{}{
			"id":           link.ID,
			"display_name": link.DisplayName,
			"link":         link.Link,
			"icon_path":    link.IconPath,
		}); err != nil {
			logrus.WithError(err).Error("unable to create link nodes")
			return errors.New("unable to create link node")
		}
		for _, tag := range link.Tags {
			if _, err := sess.Run("MATCH (a:Link),(b:Tag) WHERE a.id = $linkID AND b.id = $tagID CREATE (a)-[r:HAS_TAG]->(b)", map[string]interface{}{
				"linkID": link.ID,
				"tagID":  tag,
			}); err != nil {
				logrus.WithError(err).Error("unable to create link tag edge")
				return errors.New("unable to create link tag edge")
			}
		}
	}
	return nil
}
