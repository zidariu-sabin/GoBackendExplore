package main

import (
	app "GoBackendExploreMovieTracker/internal"
	"GoBackendExploreMovieTracker/internal/routes"
	"flag"
	"net/http"
	"time"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "go backend port")
	flag.Parse()
	app, err := app.NewApplication()

	if err != nil {
		panic(err)
	}

	defer app.DB.Close()

	app.Logger.Printf("App is successfully running \n")

	router := routes.SetupRoutes(app)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("App is successfully running on port %d \n", port)

	err = server.ListenAndServe()

	if err != nil {
		server.ErrorLog.Fatal()
		app.Logger.Fatal()
	}
	// store := store.NewPostgresSeedStore(app.DB)
}
