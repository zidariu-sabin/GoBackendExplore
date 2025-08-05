package api

import (
	"GoBackendExploreMovieTracker/internal/store"
	"GoBackendExploreMovieTracker/internal/utils"
	"GoBackendExploreMovieTracker/internal/utils/errors"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type MovieRequest struct {
	ID          *int       `json:"id"`
	Title       *string    `json:"title"`
	ReleaseDate *time.Time `json:"release_date"`
	Overview    *string    `json:"overview"`
	PosterPath  *string    `json:"poster_path"`
	GenreIds    []int      `json:"genre_ids"`
}

type MovieHandler struct {
	movieStore store.MovieStore
	logger     *log.Logger
}

func NewMovieHandler(movieStore store.MovieStore, logger *log.Logger) *MovieHandler {
	return &MovieHandler{
		movieStore: movieStore,
		logger:     logger,
	}
}

func (h *MovieHandler) HandleCreateMovie(w http.ResponseWriter, r *http.Request) {

	var movie store.Movie

	err := json.NewDecoder(r.Body).Decode(&movie)

	if err != nil {
		h.logger.Printf("ERROR: decodingCreateMovie: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, errors.ERROR_STATUS_BAD_REQUEST)
		return
	}

	err = h.movieStore.CreateMovie(&movie)

	if err != nil {
		h.logger.Println("ERROR: creatingMovie:", err)
		utils.WriteJson(w, http.StatusInternalServerError, errors.ERROR_STATUS_INTERNAL_SERVER_ERROR)
		return
	}

	utils.WriteJson(w, http.StatusCreated, utils.Envelope{"movie": movie})
}

func (h *MovieHandler) HandleGetMovieById(w http.ResponseWriter, r *http.Request) {

	id, err := utils.ReadIDParam(r)

	if err != nil {
		h.logger.Printf("ERROR: gettingMovieById: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, errors.ERROR_STATUS_BAD_REQUEST)
		return
	}

	movie, err := h.movieStore.GetMovieById(id)

	if err != nil {
		h.logger.Printf("ERROR: gettingMovieById: %v", err)
		utils.WriteJson(w, http.StatusNotFound, errors.ERROR_STATUS_NOT_FOUND)
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"movie": movie})

}
