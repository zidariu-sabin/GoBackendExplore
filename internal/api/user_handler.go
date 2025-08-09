package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"GoBackendExploreMovieTracker/internal/store"
	"GoBackendExploreMovieTracker/internal/utils"
	"GoBackendExploreMovieTracker/internal/utils/httpErrors"
)

// struct used for null field checking in registering
type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (uh *UserHandler) validateRegisterRequest(req *registerUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}

	if len(req.Username) > 50 {
		return errors.New("Username cannot be greater than 50 characters")
	}

	if req.Email == "" {
		return errors.New("Email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	return nil
}

func (uh *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {

	var req registerUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		uh.logger.Printf("ERROR: decodingRegisterRequest: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, httpErrors.ERROR_STATUS_BAD_REQUEST)
		return
	}

	err = uh.validateRegisterRequest(&req)

	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
	}

	//password logic
	err = user.PasswordHash.Set(req.Password)

	if err != nil {
		uh.logger.Printf("ERROR: hashingPassword: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, httpErrors.ERROR_STATUS_INTERNAL_SERVER_ERROR)
		return
	}

	err = uh.userStore.CreateUser(user)

	if err != nil {
		uh.logger.Printf("ERROR: registeringUser: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, httpErrors.ERROR_STATUS_INTERNAL_SERVER_ERROR)
		return
	}

	utils.WriteJson(w, http.StatusCreated, utils.Envelope{"user": user})
}
