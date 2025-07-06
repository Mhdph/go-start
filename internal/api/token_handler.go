package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mhdph/go-start/internal/store"
	"github.com/mhdph/go-start/internal/store/tokens"
	"github.com/mhdph/go-start/internal/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type crateTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req crateTokenRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userStore.GetUserByUsername(req.Username)

	if err != nil {
		h.logger.Printf("Error getting user by username: %v", err)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	passwordsDoMatch, err := user.Password.Matches(req.Password)

	if err != nil {
		h.logger.Printf("Error checking password for user: %s", req.Username)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if !passwordsDoMatch {
		h.logger.Printf("Invalid password for user: %s", req.Username)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := h.tokenStore.Create(user.ID, tokens.ScopeAuthentication, 24*time.Hour)

	if err != nil {
		h.logger.Printf("Error creating token: %v", err)
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	utils.WriteJson(w, http.StatusCreated, utils.Envelope{"auth_token": token})

}
