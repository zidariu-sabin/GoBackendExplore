package app

import (
	"GoBackendExploreMovieTracker/internal/api"
	"GoBackendExploreMovieTracker/internal/store"
	"database/sql"
	"log"
	"os"
)

type Application struct {
	Logger       *log.Logger
	DB           *sql.DB
	MovieHandler *api.MovieHandler
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()

	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	// loggerAdapter := zerologadapter.New(zerolog.New(os.Stdout))
	// db = sqldblogger.OpenDriver(dsn, db.Driver(), loggerAdapter/*, using_default_options*/)

	//stores
	movieStore := store.NewPostgresMovieStore(pgDB)

	//handlers
	movieHandler := api.NewMovieHandler(movieStore, logger)

	app := &Application{
		Logger:       logger,
		DB:           pgDB,
		MovieHandler: movieHandler,
	}

	return app, nil
}
