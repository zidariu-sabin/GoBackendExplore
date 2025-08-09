package app

import (
	"GoBackendExploreMovieTracker/internal/api"
	"GoBackendExploreMovieTracker/internal/middleware"
	"GoBackendExploreMovieTracker/internal/store"
	"database/sql"
	"log"
	"os"
)

type Application struct {
	Logger           *log.Logger
	DB               *sql.DB
	MovieHandler     *api.MovieHandler
	UserHandler      *api.UserHandler
	Tokenhandler     *api.TokenHandler
	WatchlistHandler *api.WatchlistHandler
	Middleware       *middleware.UserMiddleware
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
	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPostgrestokenStore(pgDB)
	watchlistStore := store.NewPostgresWatchlistStore(pgDB)

	//handlers
	movieHandler := api.NewMovieHandler(movieStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)
	watchlistHandler := api.NewWatchlistHandler(watchlistStore, logger)

	//middleware
	middleware := middleware.NewUserMiddleware(userStore)

	app := &Application{
		Logger:           logger,
		DB:               pgDB,
		MovieHandler:     movieHandler,
		UserHandler:      userHandler,
		Tokenhandler:     tokenHandler,
		WatchlistHandler: watchlistHandler,
		Middleware:       middleware,
	}

	return app, nil
}
