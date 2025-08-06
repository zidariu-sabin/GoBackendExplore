package routes

import (
	app "GoBackendExploreMovieTracker/internal"

	"github.com/gorilla/mux"
)

func SetupRoutes(app *app.Application) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/movie/{id}", app.MovieHandler.HandleGetMovieById).Methods("GET")
	r.HandleFunc("/movie", app.MovieHandler.HandleCreateMovie).Methods("POST")
	r.HandleFunc("/movie/{id}", app.MovieHandler.HandleUpdateMovie).Methods("PUT")
	r.HandleFunc("/movie/{id}", app.MovieHandler.HandleDeleteMovie).Methods("DELETE")

	// Define your routes here
	// Example: router.HandleFunc("/api/movies", movieHandler).Methods("GET")

	return r
}
