package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"net/http"
)

func main() {
	app := NewApp()
	app.Run()
	logrus.WithField("error", <-app.Done()).Error("terminated")
}

func NewApp() *fx.App {

	config := LoadConfig()
	providers := []interface{}{
		config.LoadBucketConfig,
		NewBookService,
		NewBucketBookStorage,
		NewMux,
	}

	//TODO db type enum
	switch config.DatabaseType {
	case "sqlite":
		providers = append(providers, config.LoadSQLiteDatabaseConfig)
		providers = append(providers, NewSQLiteDatabase)
	case "postgres":
		providers = append(providers, config.LoadPostgresDatabaseConfig)
		providers = append(providers, NewPostgresDatabase)
		providers = append(providers, NewPostgresBookRepository)
	}
	return fx.New(
		fx.Provide(
			providers...,
		),
		fx.Invoke(MakeBookHandler),
		fx.Logger(NewLogger()),
	)
}
func NewMux(lc fx.Lifecycle) *mux.Router {
	logrus.Info("creating mux")

	mux := mux.NewRouter()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logrus.Info("starting server")
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logrus.Info("stopping server")
			return server.Shutdown(ctx)
		},
	})

	return mux
}

func NewLogger() *logrus.Logger {
	return logrus.New()
}
