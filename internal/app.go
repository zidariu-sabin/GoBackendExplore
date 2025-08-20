package app

import (
	"GoBackendExploreMovieTracker/internal/api"
	"GoBackendExploreMovieTracker/internal/crons"
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
	TokenHandler     *api.TokenHandler
	WatchlistHandler *api.WatchlistHandler
	ReviewHandler    *api.ReviewHandler
	RatingHandler    *api.RatingHandler
	Middleware       *middleware.UserMiddleware
	CronJobPipeline  *crons.CronJobPipeLine
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
	tokenStore := store.NewPostgresTokenStore(pgDB)
	watchlistStore := store.NewPostgresWatchlistStore(pgDB)
	reviewStore := store.NewPostgresReviewStore(pgDB)
	ratingStore := store.NewPostgresRatingStore(pgDB)

	//handlers
	movieHandler := api.NewMovieHandler(movieStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)
	watchlistHandler := api.NewWatchlistHandler(watchlistStore, logger)
	reviewHandler := api.NewReviewHandler(reviewStore, logger)
	ratingHandler := api.NewRatingHandler(ratingStore, logger)

	//middleware
	middleware := middleware.NewUserMiddleware(userStore)

	//cron jobs
	cronJobStore := store.NewPostgresCronJobStore(pgDB)
	cronJobPipeline := crons.NewCronJobPipeline(crons.RunningCrons, cronJobStore, logger)

	app := &Application{
		Logger:           logger,
		DB:               pgDB,
		MovieHandler:     movieHandler,
		UserHandler:      userHandler,
		TokenHandler:     tokenHandler,
		WatchlistHandler: watchlistHandler,
		ReviewHandler:    reviewHandler,
		RatingHandler:    ratingHandler,
		Middleware:       middleware,
		CronJobPipeline:  cronJobPipeline,
	}

	return app, nil
}
