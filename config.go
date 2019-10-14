package main

import "os"

type Config struct {
	DatabaseType string
	SQLiteConfig SQLiteDatabaseConfig
}

func NewConfig() Config {
	dbType := getEnv("DB_TYPE", "sqlite")

	config := &Config{
		DatabaseType: getEnv("DB_TYPE", "sqlite"),
	}

	switch dbType {
	case "sqlite":
		config.SQLiteConfig = NewSQLiteDatabaseConfig()
	}

	return *config
}

type SQLiteDatabaseConfig struct {
	File string
}

func NewSQLiteDatabaseConfig() SQLiteDatabaseConfig {
	//TODO check file, maybe init
	return SQLiteDatabaseConfig{
		File: getEnv("DB_FILE", "books"),
	}
}

func getEnv(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
