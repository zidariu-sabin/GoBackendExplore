package main

import (
	app "GoBackendExploreMovieTracker/internal"
)

func main() {

	app, err := app.NewApplication()

	if err != nil {
		panic(err)
	}

	defer app.DB.Close()

	app.Logger.Printf("App is successfully running \n")

	// store := store.NewPostgresWorkoutStore(app.DB)
}
