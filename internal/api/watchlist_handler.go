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

type WatchlistEntryRequest struct {
	MovieId *int64 `json:"movie_id"`
}

type WatchlistHandler struct {
	watchlistStore *store.PostgresWatchlistStore
	logger         *log.Logger
}

func NewWatchlistHandler(watchlistStore *store.PostgresWatchlistStore, logger *log.Logger) *WatchlistHandler {
	return &WatchlistHandler{
		watchlistStore: watchlistStore,
		logger:         logger,
	}
}

func (h WatchlistHandler) HandleAddToWatchlist(w http.ResponseWriter, r *http.Request) {

	var req WatchlistEntryRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		h.logger.Printf("ERROR: decodingAddToWatchlist: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, httpErrors.ERROR_STATUS_BAD_REQUEST)
		return
	}

	user := middleware.GetUser(r)

	watchlistEntry := &store.WatchlistEntry{
		UserID:  user.ID,
		MovieID: *req.MovieId,
	}

	err = h.watchlistStore.AddToWatchlist(watchlistEntry.UserID, watchlistEntry.MovieID)

	if err != nil {
		h.logger.Printf("ERROR: addingToWatchlist: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, httpErrors.ERROR_STATUS_INTERNAL_SERVER_ERROR)
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"movie": "added to watchlist"})

}

func (h WatchlistHandler) HandleRemoveFromWatchlist(w http.ResponseWriter, r *http.Request) {

	var req WatchlistEntryRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		h.logger.Printf("ERROR: decodingRemoveFromWatchlist: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, httpErrors.ERROR_STATUS_BAD_REQUEST)
		return
	}
	user := middleware.GetUser(r)

	watchlistEntry := &store.WatchlistEntry{
		UserID:  user.ID,
		MovieID: *req.MovieId,
	}

	err = h.watchlistStore.RemoveFromWatchlist(watchlistEntry.UserID, watchlistEntry.MovieID)

	if err != nil {
		h.logger.Printf("ERROR: removingFromWatchList: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, httpErrors.ERROR_STATUS_INTERNAL_SERVER_ERROR)
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"movie": "removed from watchlist"})

}

func (h WatchlistHandler) HandleGetWatchlist(w http.ResponseWriter, r *http.Request) {

	userID, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: retrievingIdUserWatchlist: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, httpErrors.ERROR_STATUS_BAD_REQUEST)
		return
	}

	movies, err := h.watchlistStore.GetWatchlist(userID)

	if err != nil {
		h.logger.Printf("ERROR: gettingWatchList: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, httpErrors.ERROR_STATUS_INTERNAL_SERVER_ERROR)
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"movies": movies})

}
