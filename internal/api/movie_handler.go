package api

import (
	"GoBackendExploreMovieTracker/internal/store"
	"GoBackendExploreMovieTracker/internal/utils"
	"GoBackendExploreMovieTracker/internal/utils/httpErrors"
	"database/sql"
	"encoding/json"
	"errors"
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
	GenreIds    []int64    `json:"genre_ids"`
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
		utils.WriteJson(w, http.StatusBadRequest, httpErrors.ERROR_STATUS_BAD_REQUEST)
		return
	}

	err = h.movieStore.CreateMovie(&movie)

	if err != nil {
		h.logger.Println("ERROR: creatingMovie:", err)
		utils.WriteJson(w, http.StatusInternalServerError, httpErrors.ERROR_STATUS_INTERNAL_SERVER_ERROR)
		return
	}

	utils.WriteJson(w, http.StatusCreated, utils.Envelope{"movie": movie})
}

func (h *MovieHandler) HandleGetMovieById(w http.ResponseWriter, r *http.Request) {

	id, err := utils.ReadIDParam(r)

	if err != nil {
		h.logger.Printf("ERROR: gettingMovieById: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, httpErrors.ERROR_STATUS_BAD_REQUEST)
		return
	}

	movie, err := h.movieStore.GetMovieById(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.logger.Printf("WARN: movie not found by id: %v", err)
			utils.WriteJson(w, http.StatusNotFound, httpErrors.ERROR_STATUS_NOT_FOUND)
			return
		} else {
			h.logger.Printf("ERROR: gettingMovieById: %v", err)
			utils.WriteJson(w, http.StatusInternalServerError, httpErrors.ERROR_STATUS_INTERNAL_SERVER_ERROR)
			return
		}
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"movie": movie})

}

func (h *MovieHandler) HandleUpdateMovie(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParam(r)

	if err != nil {
		h.logger.Printf("ERROR: updatingMovie: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, httpErrors.ERROR_STATUS_BAD_REQUEST)
		return
	}

	existingMovie, err := h.movieStore.GetMovieById(id)

	if errors.Is(err, sql.ErrNoRows) {
		h.logger.Printf("WARN: movie to update not found: %v", err)
		utils.WriteJson(w, http.StatusNotFound, httpErrors.ERROR_STATUS_NOT_FOUND)
	} else {
		h.logger.Printf("ERROR: updatingMovieGettingById: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, httpErrors.ERROR_STATUS_INTERNAL_SERVER_ERROR)
	}

	var moviePayload MovieRequest
	err = json.NewDecoder(r.Body).Decode(&moviePayload)

	if moviePayload.Title != nil {
		existingMovie.Title = *moviePayload.Title
	}

	if moviePayload.ReleaseDate != nil {
		existingMovie.ReleaseDate = *moviePayload.ReleaseDate
	}

	if moviePayload.Overview != nil {
		existingMovie.Overview = *moviePayload.Overview
	}

	if moviePayload.PosterPath != nil {
		existingMovie.PosterPath = *moviePayload.PosterPath
	}

	if moviePayload.GenreIds != nil {
		existingMovie.GenreIds = moviePayload.GenreIds
	}

	if err != nil {
		h.logger.Printf("ERROR: decodingUpdateMovie: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, httpErrors.ERROR_STATUS_BAD_REQUEST)
		return
	}

	err = h.movieStore.UpdateMovie(existingMovie)

	if err != nil {
		h.logger.Printf("ERROR: updatingMovie: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, httpErrors.ERROR_STATUS_INTERNAL_SERVER_ERROR)
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"movie": existingMovie})
}

func (h *MovieHandler) HandleDeleteMovie(w http.ResponseWriter, r *http.Request) {

	id, err := utils.ReadIDParam(r)

	if err != nil {
		h.logger.Printf("ERROR: gettingMovieById: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, httpErrors.ERROR_STATUS_BAD_REQUEST)
		return
	}

	err = h.movieStore.DeleteMovie(id)

	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			h.logger.Printf("ERROR: gettingMovieById: %v", err)
			utils.WriteJson(w, http.StatusNotFound, httpErrors.ERROR_STATUS_NOT_FOUND)
		default:
			h.logger.Printf("ERROR: gettingMovieById: %v", err)
			utils.WriteJson(w, http.StatusInternalServerError, httpErrors.ERROR_STATUS_INTERNAL_SERVER_ERROR)
		}
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"message": "movie deleted successfully"})

}
