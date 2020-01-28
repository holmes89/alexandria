package main

import (
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
	return fx.New(
		fx.Provide(
			NewBookBucket,
			NewSiteBucket,
			NewService,
			NewMux,
		),
		fx.Invoke(MakeThumnbnailHandler),
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