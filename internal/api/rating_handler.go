package api

import (
	"GoBackendExploreMovieTracker/internal/middleware"
	"GoBackendExploreMovieTracker/internal/store"
	"GoBackendExploreMovieTracker/internal/utils"
	"GoBackendExploreMovieTracker/internal/utils/httpErrors"
	"encoding/json"
	"log"
	"net/http"
)

type RatingRequest struct {
	UserID  *int64 `json:"user_id"`
	MovieID *int64 `json:"movie_id"`
	Score   *int64 `json:"score"`
}

type RatingHandler struct {
	ratingStore store.RatingStore
	logger      *log.Logger
}

func NewRatingHandler(ratingStore store.RatingStore, logger *log.Logger) *RatingHandler {
	return &RatingHandler{
		ratingStore: ratingStore,
		logger:      logger,
	}
}

func (h *RatingHandler) HandleCreateRating(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	var req RatingRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: decodingCreateReview: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, httpErrors.ERROR_STATUS_BAD_REQUEST)
		return
	}

	rating := &store.Rating{
		UserID:  user.ID,
		MovieID: *req.MovieID,
		Score:   *req.Score,
	}

	err = h.ratingStore.CreateRating(rating)
	if err != nil {
		h.logger.Printf("ERROR: creatingRating: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, httpErrors.ERROR_STATUS_INTERNAL_SERVER_ERROR)
		return
	}

	utils.WriteJson(w, http.StatusCreated, utils.Envelope{"rating": rating})
}

func (h *RatingHandler) HandleGetMovieRatingScore(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: updatingReview: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, httpErrors.ERROR_STATUS_BAD_REQUEST)
		return
	}

	score, ratingNo, err := h.ratingStore.GetMovieRatingScore(id)
	if err != nil {
		h.logger.Printf("ERROR: gettingMovieRatingScore: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, httpErrors.ERROR_STATUS_INTERNAL_SERVER_ERROR)
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"score": score, "rating_count": ratingNo})
}

func (h *RatingHandler) HandleUpdateRating(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	var req RatingRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		h.logger.Printf("ERROR: decodingUpdateRating: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, httpErrors.ERROR_STATUS_BAD_REQUEST)
		return
	}

	rating := &store.Rating{
		UserID:  user.ID,
		MovieID: *req.MovieID,
		Score:   *req.Score,
	}

	err = h.ratingStore.UpdateRating(rating)
	if err != nil {
		h.logger.Printf("ERROR: updatingRating: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, httpErrors.ERROR_STATUS_INTERNAL_SERVER_ERROR)
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"rating": rating})
}

func (h *RatingHandler) HandleDeleteRating(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	var req RatingRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		h.logger.Printf("ERROR: decodingDeleteRating: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, httpErrors.ERROR_STATUS_BAD_REQUEST)
		return
	}

	err = h.ratingStore.DeleteRating(user.ID, *req.MovieID)
	if err != nil {
		h.logger.Printf("ERROR: deletingRating: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, httpErrors.ERROR_STATUS_INTERNAL_SERVER_ERROR)
		return
	}

	utils.WriteJson(w, http.StatusNoContent, nil)
}
