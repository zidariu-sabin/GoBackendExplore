package routes

import (
	app "GoBackendExploreMovieTracker/internal"

	"github.com/gorilla/mux"
)

func SetupRoutes(app *app.Application) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/user", app.UserHandler.HandleRegisterUser).Methods("POST")
	r.HandleFunc("/tokens/authentication", app.TokenHandler.HandleCreateToken).Methods("POST")

	r.Use(app.Middleware.Authenticate)
	r.HandleFunc("/movie/{id}", app.Middleware.RequireUser(app.MovieHandler.HandleGetMovieById)).Methods("GET")
	r.HandleFunc("/movie", app.Middleware.RequireUser(app.MovieHandler.HandleCreateMovie)).Methods("POST")
	r.HandleFunc("/movie/{id}", app.Middleware.RequireUser(app.MovieHandler.HandleUpdateMovie)).Methods("PUT")
	r.HandleFunc("/movie/{id}", app.Middleware.RequireUser(app.MovieHandler.HandleDeleteMovie)).Methods("DELETE")

	r.HandleFunc("/watchlist", app.Middleware.RequireUser(app.WatchlistHandler.HandleAddToWatchlist)).Methods("POST")
	r.HandleFunc("/watchlist", app.Middleware.RequireUser(app.WatchlistHandler.HandleRemoveFromWatchlist)).Methods("DELETE")
	r.HandleFunc("/watchlist/{id}", app.Middleware.RequireUser(app.WatchlistHandler.HandleGetWatchlist)).Methods("GET")

	r.HandleFunc("/review", app.Middleware.RequireUser(app.ReviewHandler.HandleCreateReview)).Methods("POST")
	r.HandleFunc("/review/{id}", app.Middleware.RequireUser(app.ReviewHandler.HandleGetReviewById)).Methods("GET")
	r.HandleFunc("/review/{id}", app.Middleware.RequireUser(app.ReviewHandler.HandleUpdateReview)).Methods("PUT")
	r.HandleFunc("/review/{id}", app.Middleware.RequireUser(app.ReviewHandler.HandleDeleteReview)).Methods("DELETE")

	r.HandleFunc("/rating", app.Middleware.RequireUser(app.RatingHandler.HandleCreateRating)).Methods("POST")
	r.HandleFunc("/rating/movie/{id}", app.Middleware.RequireUser(app.RatingHandler.HandleGetMovieRatingScore)).Methods("GET")
	r.HandleFunc("/rating", app.Middleware.RequireUser(app.RatingHandler.HandleUpdateRating)).Methods("PUT")
	r.HandleFunc("/rating", app.Middleware.RequireUser(app.RatingHandler.HandleDeleteRating)).Methods("DELETE")
	return r
}
