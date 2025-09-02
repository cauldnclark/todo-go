package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cauldnclark/todo-go/internal/middleware"
	"github.com/cauldnclark/todo-go/internal/models"
	"github.com/cauldnclark/todo-go/internal/service"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	userService *service.UserService
	validator   *validator.Validate
}

func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		validator:   validator.New(),
	}
}

func (h *AuthHandler) GoogleSignIn(w http.ResponseWriter, r *http.Request) {
	var req models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		http.Error(w, "Validation failed"+err.Error(), http.StatusBadRequest)
		return
	}

	authResp, err := h.userService.AuthenticateWithGoogle(r.Context(), req.Token)
	if err != nil {
		http.Error(w, "Authentication failed"+err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(authResp); err != nil {
		http.Error(w, "Failed to encode response"+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get user"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode response"+err.Error(), http.StatusInternalServerError)
		return
	}
}
