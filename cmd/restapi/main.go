package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"bitbucket.org/faceit/internal/events/awssns"

	"bitbucket.org/faceit/app"
	"bitbucket.org/faceit/internal/users"
	"github.com/aws/aws-sdk-go/aws/session"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/subosito/gotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const defaultPORT = "8080"

func init() {
	gotenv.Load()
	app.Init(os.Getenv("APP_ENV"), os.Getenv("APP_URL"))
}

func main() {
	db := initGORM()

	awsSession := session.Must(session.NewSession())

	// services
	var eventService = awssns.New(awsSession, mustGetEnv("EVENT_TOPIC"))

	var userervice *users.Service
	{
		repo := users.NewMysqlRepo(db)
		userervice = &users.Service{Users: repo}
	}

	r := mux.NewRouter()

	r.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	// handlers

	// users handler
	{
		h := users.Handler{
			Users:  userervice,
			Events: eventService,
		}
		r.HandleFunc("/users", h.List).Methods(http.MethodGet)
		r.HandleFunc("/users", h.Create).Methods(http.MethodPost)
		r.HandleFunc("/users/{id}", h.Update).Methods(http.MethodPut)
		r.HandleFunc("/users/{id}", h.Delete).Methods(http.MethodDelete)

	}

	// start the http server
	log.Printf("Listening on port: %s", getADDR())

	sentryHandler := sentryhttp.New(sentryhttp.Options{Repanic: true})

	handler := cors.AllowAll().Handler(sentryHandler.Handle(r))

	if err := http.ListenAndServe(getADDR(), handler); err != nil {
		log.Fatal(err)
	}
}

func getADDR() string {
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return ":" + defaultPORT
}

func mustGetEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	panic("env required: " + key)
}

func mysqlDSN(dsn string) string {
	u, err := url.Parse(dsn)
	if err != nil {
		log.Fatalf("mysql dsn parse error: %s", err)
	}
	userPass := u.User.Username()
	if pass, ok := u.User.Password(); ok {
		userPass += ":" + pass
	}
	return fmt.Sprintf("%s@tcp(%s)%s?charset=utf8&parseTime=True&loc=Local", userPass, u.Host, u.Path)
}

func initGORM() *gorm.DB {
	logmode := logger.Info
	if app.ENV == app.PROD {
		logmode = logger.Silent
	}

	db, err := gorm.Open(mysql.Open(mysqlDSN(mustGetEnv("MYSQL_DSN"))), &gorm.Config{
		Logger: logger.Default.LogMode(logmode),
	})
	if err != nil {
		log.Fatalf("gorm: %s", err)
	}

	return db
}
