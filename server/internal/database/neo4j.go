package database

import (
	"alexandria/internal/common"
	"context"
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
		return neo4j.NewDriver(config.URI, neo4j.BasicAuth(config.Username, config.Password, ""))
	})
	if err != nil {
		logrus.WithError(err).Fatal("unable to connect to postgres")
	}
	logrus.Info("connected to postgres")
	db := &Neo4jDatabase{driver}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logrus.Info("closing connection for neo4j")
			driver.Close()
			return nil
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

		logrus.WithError(err).Error("error connecting to postgres, retrying")
	}
	return nil, fmt.Errorf("after %d attempts, connection failed", attempts)
}
