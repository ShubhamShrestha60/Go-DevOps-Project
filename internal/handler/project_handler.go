package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/user/devpulse/internal/middleware"
	"github.com/user/devpulse/internal/models"
	"github.com/user/devpulse/internal/service"
)

type ProjectHandler struct {
	service *service.ProjectService
}

func NewProjectHandler(s *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: s}
}

// Create godoc
// @Summary Create a new project
// @Description Create a new project for the current user
// @Tags projects
// @Accept json
// @Produce json
// @Param project body object true "Project info"
// @Success 201 {object} models.Project
// @Failure 400 {string} string "Bad Request"
// @Router /api/projects [post]
func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserID(r.Context())
	project, err := h.service.CreateProject(r.Context(), input.Name, input.Description, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(project)
}

// List godoc
// @Summary List projects
// @Description List projects for the current user or search by query
// @Tags projects
// @Produce json
// @Param q query string false "Search query"
// @Success 200 {array} models.Project
// @Router /api/projects [get]
func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	var projects []*models.Project
	var err error

	if query != "" {
		projects, err = h.service.SearchProjects(r.Context(), query)
	} else {
		userID := middleware.GetUserID(r.Context())
		projects, err = h.service.ListUserProjects(r.Context(), userID)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(projects)
}

// Get godoc
// @Summary Get project details
// @Description Get a single project by ID
// @Tags projects
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} models.Project
// @Failure 404 {string} string "Not Found"
// @Router /api/projects/{id} [get]
func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	project, err := h.service.GetProject(r.Context(), id)
	if err != nil {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(project)
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserID(r.Context())
	if err := h.service.UpdateProject(r.Context(), userID, id, input.Name, input.Description); err != nil {
		if err == models.ErrUnauthorized {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserID(r.Context())
	if err := h.service.DeleteProject(r.Context(), userID, id); err != nil {
		if err == models.ErrUnauthorized {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Stats godoc
// @Summary Get project stats
// @Description Get count of projects
// @Tags projects
// @Produce json
// @Success 200 {object} models.ProjectStats
// @Router /api/projects/stats [get]
func (h *ProjectHandler) Stats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(stats)
}
