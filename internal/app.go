package app

import (
	"GoBackendExploreMovieTracker/internal/store"
	"database/sql"
	"log"
	"os"
)

type Application struct {
	Logger *log.Logger
	DB     *sql.DB
}

func NewApplication() (*Application, error) {
	pgDb, err := store.Open()

	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	// loggerAdapter := zerologadapter.New(zerolog.New(os.Stdout))
	// db = sqldblogger.OpenDriver(dsn, db.Driver(), loggerAdapter/*, using_default_options*/)

	app := &Application{
		Logger: logger,
		DB:     pgDb,
	}

	return app, nil
}
