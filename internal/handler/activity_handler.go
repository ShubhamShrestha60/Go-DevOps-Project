package handler

import (
	"encoding/json"
	"net/http"

	"github.com/user/devpulse/internal/service"
)

type ActivityHandler struct {
	service *service.ActivityService
}

func NewActivityHandler(s *service.ActivityService) *ActivityHandler {
	return &ActivityHandler{service: s}
}

func (h *ActivityHandler) List(w http.ResponseWriter, r *http.Request) {
	activities, err := h.service.GetRecent(r.Context(), 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activities)
}
