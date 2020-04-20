package main

import (
	"alexandria/internal/backup"
	"alexandria/internal/books"
	"alexandria/internal/common"
	"alexandria/internal/database"
	"alexandria/internal/documents"
	"alexandria/internal/journal"
	"alexandria/internal/links"
	"alexandria/internal/papers"
	"alexandria/internal/tags"
	"alexandria/internal/user"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"net/http"
	"os"
	"strings"
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
			config.LoadNeo4jConfig,
			common.NewGCPBucketStorage,
			common.NewBucketDocumentStorage,
			common.NewBackupStorage,
			documents.NewDocumentService,
			books.NewBookService,
			papers.NewPaperService,
			links.NewService,
			backup.NewService,
			database.NewPostgresDatabase,
			database.NewNeo4jDatabase,
			database.NewPostgresDocumentRepository,
			database.NewUserPostgresRepository,
			database.NewJournalRepository,
			database.NewLinksRepository,
			database.NewTagsRepository,
			database.NewBackupRepository,
			user.NewUserService,
			NewMux,
		),
		fx.Invoke(documents.MakeDocumentHandler,
			books.MakeBookHandler,
			user.MakeLoginHandler,
			papers.MakePaperHandler,
			journal.MakeJournalHandler,
			links.MakeLinksHandler,
			tags.MakeLinksHandler,
			backup.MakeBackupHandler,
			backup.NewBackupRunner,
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

	router.Use(cors, authenticate)
	handler := (cors)((authenticate)(router))

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

func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Exclude auth
		if strings.Contains(r.URL.Path, "auth") && r.Method == "GET" {
			next.ServeHTTP(w, r) // call original
			return
		}

		// sample token string taken from the New example
		tokenString := r.Header.Get("Authorization")
		tokenString = strings.Replace(tokenString, "Bearer ", "", -1)
		if tokenString == "" {
			http.Error(w, "Authorization Header Required", http.StatusUnauthorized)
			return
		}
		// Parse takes the token string and a function for looking up the key. The latter is especially
		// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
		// head of the token to identify which key to use, but the parsed token (head and claims) is provided
		// to the callback, providing flexibility.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			next.ServeHTTP(w, r) // call original
		} else {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}
	})
}
