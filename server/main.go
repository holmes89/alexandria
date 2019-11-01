package main

import (
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
		NewBucketStorage,
		NewBucketBookStorage,
		NewBookService,
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

	//secureMiddleware := secure.New(secure.Options{
	//	AllowedHosts:          []string{"localhost:8080"},
	//	AllowedHostsAreRegex:  true,
	//	SSLRedirect:           false,
	//	STSSeconds:            31536000,
	//})
	router := mux.NewRouter()

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logrus.Info("starting server")
			go http.ListenAndServe(":8080", handlers.CORS()(router))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logrus.Info("stopping server")
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT)
			return fmt.Errorf("%s", <-c)
		},
	})

	return router
}

func NewLogger() *logrus.Logger {
	return logrus.New()
}
