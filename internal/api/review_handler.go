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
)

type ReviewRequest struct {
	ID      *int64  `json:"id"`
	UserId  *int64  `json:"user_id"`
	MovieId *int64  `json:"movie_id"`
	Content *string `json:"content"`
}

type ReviewHandler struct {
	reviewStore store.ReviewStore
	logger      *log.Logger
}

func NewReviewHandler(reviewStore store.ReviewStore, logger *log.Logger) *ReviewHandler {
	return &ReviewHandler{
		reviewStore: reviewStore,
		logger:      logger,
	}
}

func (h *ReviewHandler) HandleCreateReview(w http.ResponseWriter, r *http.Request) {
	var reviewPayload ReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&reviewPayload); err != nil {
		h.logger.Printf("ERROR: decodingCreateReview: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, httpErrors.ERROR_STATUS_BAD_REQUEST)
		return
	}

	review := &store.Review{}
	if reviewPayload.UserId != nil {
		review.UserId = *reviewPayload.UserId
	}
	if reviewPayload.MovieId != nil {
		review.MovieId = *reviewPayload.MovieId
	}
	if reviewPayload.Content != nil {
		review.Content = *reviewPayload.Content
	}

	if err := h.reviewStore.CreateReview(review); err != nil {
		h.logger.Printf("ERROR: creatingReview: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, httpErrors.ERROR_STATUS_INTERNAL_SERVER_ERROR)
		return
	}

	utils.WriteJson(w, http.StatusCreated, utils.Envelope{"review": review})
}

func (h *ReviewHandler) HandleGetReviewById(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: gettingReviewById: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, httpErrors.ERROR_STATUS_BAD_REQUEST)
		return
	}

	review, err := h.reviewStore.GetReviewById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.logger.Printf("WARN: review not found by id: %v", err)
			utils.WriteJson(w, http.StatusNotFound, httpErrors.ERROR_STATUS_NOT_FOUND)
			return
		} else {
			h.logger.Printf("ERROR: gettingReviewById: %v", err)
			utils.WriteJson(w, http.StatusInternalServerError, httpErrors.ERROR_STATUS_INTERNAL_SERVER_ERROR)
			return
		}
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"review": review})
}

func (h *ReviewHandler) HandleUpdateReview(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: updatingReview: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, httpErrors.ERROR_STATUS_BAD_REQUEST)
		return
	}

	existingReview, err := h.reviewStore.GetReviewById(id)
	if errors.Is(err, sql.ErrNoRows) {
		h.logger.Printf("WARN: review to update not found: %v", err)
		utils.WriteJson(w, http.StatusNotFound, httpErrors.ERROR_STATUS_NOT_FOUND)
		return
	} else if err != nil {
		h.logger.Printf("ERROR: updatingReviewGettingById: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, httpErrors.ERROR_STATUS_INTERNAL_SERVER_ERROR)
		return
	}

	var reviewPayload ReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&reviewPayload); err != nil {
		h.logger.Printf("ERROR: decodingUpdateReview: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, httpErrors.ERROR_STATUS_BAD_REQUEST)
		return
	}

	if reviewPayload.Content != nil {
		existingReview.Content = *reviewPayload.Content
	}

	if err := h.reviewStore.UpdateReview(existingReview); err != nil {
		h.logger.Printf("ERROR: updatingReview: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, httpErrors.ERROR_STATUS_INTERNAL_SERVER_ERROR)
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"review": existingReview})
}

func (h *ReviewHandler) HandleDeleteReview(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: deletingReview: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, httpErrors.ERROR_STATUS_BAD_REQUEST)
		return
	}

	if err := h.reviewStore.DeleteReview(id); err != nil {
		switch {
		case err == sql.ErrNoRows:
			h.logger.Printf("ERROR: review not found: %v", err)
			utils.WriteJson(w, http.StatusNotFound, httpErrors.ERROR_STATUS_NOT_FOUND)
		default:
			h.logger.Printf("ERROR: deletingReview: %v", err)
			utils.WriteJson(w, http.StatusInternalServerError, httpErrors.ERROR_STATUS_INTERNAL_SERVER_ERROR)
		}
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"message": "review deleted successfully"})
}
