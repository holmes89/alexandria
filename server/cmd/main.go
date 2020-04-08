package main

import (
	"alexandria/internal/common"
	"alexandria/internal/database"
	"alexandria/internal/documents"
	"alexandria/internal/journal"
	"alexandria/internal/user"
	"context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"net/http"
)

func main() {
	app := NewApp()
	app.Run()
}

func NewApp() *fx.App {

	config := common.LoadConfig()

	return fx.New(
		fx.Provide(
			config.LoadBucketConfig,
			config.LoadPostgresDatabaseConfig,
			common.NewGCPBucketStorage,
			common.NewBucketDocumentStorage,
			documents.NewDocumentService,
			documents.NewBookService,
			documents.NewPaperService,
			database.NewPostgresDatabase,
			database.NewPostgresDocumentRepository,
			database.NewUserPostgresRepository,
			database.NewJournalRepository,
			user.NewUserService,
			NewMux,
		),
		fx.Invoke(documents.MakeDocumentHandler,
			documents.MakeBookHandler,
			user.MakeLoginHandler,
			documents.MakePaperHandler,
			journal.MakeJournalHandler,
		),
		fx.Logger(NewLogger()),
	)
}
func NewMux(lc fx.Lifecycle) *mux.Router {
	logrus.Info("creating mux")

	router := mux.NewRouter()

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "PATCH", "OPTIONS", "DELETE"})
	cors := handlers.CORS(originsOk, headersOk, methodsOk)

	router.Use(cors)
	handler := (cors)(router)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logrus.Info("starting server")
			go http.ListenAndServe(":8080", handler)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logrus.Info("stopping server")
			return nil
		},
	})

	return router
}

//NewLogger uses logrus for logging
func NewLogger() *logrus.Logger {
	return logrus.New()
}

// EndpointLogging middleware to handle logging and control headers.
func EndpointLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			return
		}
		url := r.URL.String()
		logrus.WithFields(logrus.Fields{"uri": url, "method": r.Method}).Info("endpoint")
		h.ServeHTTP(w, r)
	})
}
