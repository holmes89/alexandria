package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

type Config struct {
	DatabaseType   string
	SQLiteConfig   SQLiteDatabaseConfig
	PostgresConfig PostgresDatabaseConfig
	BucketConnectionConfig BucketConfig
}

func LoadConfig() Config {
	dbType := getEnv("DB_TYPE", "sqlite")

	config := &Config{
		DatabaseType: getEnv("DB_TYPE", "sqlite"),
	}

	logrus.WithField("type", dbType).Info("loading database config")
	return *config
}

type SQLiteDatabaseConfig struct {
	File string
}

func (c *Config) LoadSQLiteDatabaseConfig() SQLiteDatabaseConfig {
	//TODO check file, maybe init
	return SQLiteDatabaseConfig{
		File: getEnv("DB_FILE", "books"),
	}
}

type PostgresDatabaseConfig struct {
	ConnectionString string
}

func (c *Config) LoadPostgresDatabaseConfig() PostgresDatabaseConfig {

	host := getEnv("DB_HOST", "localhost")
	portString := getEnv("DB_PORT", "5432")
	port, err := strconv.Atoi(portString)
	if err != nil {
		logrus.WithError(err).Warn("invalid port falling back to default 5432")
		port = 5432
	}
	username := getEnv("DB_POSTGRES", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "postgres")
	sslmode := getEnv("DB_SSL", "disable")

	connStringFormat := "host=%s port=%d user=%s " +
		"password=%s dbname=%s sslmode=%s connect_timeout=15"

	connString := fmt.Sprintf(connStringFormat,
		host, port, username, password, dbname, sslmode)

	return PostgresDatabaseConfig{
		ConnectionString: connString,
	}
}

type BucketConfig struct {
	ConnectionString string
}

func (c *Config) LoadBucketConfig() BucketConfig {
	host := getEnv("BUCKET_HOST", "s3://my-books")
	return BucketConfig{host}
}

func getEnv(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
