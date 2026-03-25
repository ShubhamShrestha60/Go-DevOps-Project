package handler

import (
	"encoding/json"
	"net/http"

	"github.com/user/devpulse/internal/middleware"
	"github.com/user/devpulse/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// GetProfile godoc
// @Summary Get current user profile
// @Description Get the profile for the authenticated user
// @Tags users
// @Produce json
// @Success 200 {object} models.User
// @Router /api/users/profile [get]
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	id := middleware.GetUserID(r.Context())
	user, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the profile for the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Param profile body object true "Profile update info"
// @Success 204 {string} string "No Content"
// @Router /api/users/profile [put]
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	id := middleware.GetUserID(r.Context())
	var input struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateUser(r.Context(), id, input.FullName, input.Email, input.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// List godoc
// @Summary List all users
// @Description Get a list of all registered users
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Router /api/users [get]
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(users)
}

// DeleteProfile godoc
// @Summary Delete current user profile
// @Description Delete the authenticated user account
// @Tags users
// @Success 204 {string} string "No Content"
// @Router /api/users/profile [delete]
func (h *UserHandler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	id := middleware.GetUserID(r.Context())
	if err := h.service.DeleteUser(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
