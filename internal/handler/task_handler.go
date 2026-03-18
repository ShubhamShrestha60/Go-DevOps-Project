package handler

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/user/devpulse/internal/models"
	"github.com/user/devpulse/internal/service"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(s *service.TaskService) *TaskHandler {
	return &TaskHandler{service: s}
}

// Create godoc
// @Summary Create a new task
// @Description Create a new task within a project
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body object true "Task info"
// @Success 201 {object} models.Task
// @Failure 400 {string} string "Bad Request"
// @Router /api/tasks [post]
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ProjectID   string `json:"project_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    string `json:"priority"`
		AssignedTo  string `json:"assigned_to"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	projectID, err := uuid.Parse(input.ProjectID)
	if err != nil {
		http.Error(w, "invalid project_id", http.StatusBadRequest)
		return
	}

	var assignedTo *uuid.UUID
	if input.AssignedTo != "" {
		id, err := uuid.Parse(input.AssignedTo)
		if err == nil {
			assignedTo = &id
		}
	}

	task, err := h.service.CreateTask(r.Context(), projectID, input.Title, input.Description, input.Priority, assignedTo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// List godoc
// @Summary List tasks
// @Description List all tasks or filter by project, search query, or priority
// @Tags tasks
// @Produce json
// @Param project_id query string false "Project ID filter"
// @Param q query string false "Search query"
// @Param priority query string false "Priority filter"
// @Success 200 {array} models.Task
// @Router /api/tasks [get]
func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	projectIDStr := r.URL.Query().Get("project_id")
	priority := r.URL.Query().Get("priority")
	
	var tasks []*models.Task
	var err error

	if query != "" {
		tasks, err = h.service.SearchTasks(r.Context(), query)
	} else {
		var projectID *uuid.UUID
		if projectIDStr != "" {
			id, parseErr := uuid.Parse(projectIDStr)
			if parseErr != nil {
				http.Error(w, "invalid project_id", http.StatusBadRequest)
				return
			}
			projectID = &id
		}
		tasks, err = h.service.ListFiltered(r.Context(), projectID, priority)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// Get godoc
// @Summary Get task details
// @Description Get a single task by ID
// @Tags tasks
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} models.Task
// @Failure 404 {string} string "Not Found"
// @Router /api/tasks/{id} [get]
func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, err := h.service.GetTask(r.Context(), id)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(task)
}

// Update godoc
// @Summary Update task
// @Description Update task details (title, description, status, priority, assignment)
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param task body object true "Updated task info"
// @Success 204 {string} string "No Content"
// @Router /api/tasks/{id} [put]
func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		Priority    string `json:"priority"`
		AssignedTo  string `json:"assigned_to"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var assignedTo *uuid.UUID
	if input.AssignedTo != "" {
		parsedID, err := uuid.Parse(input.AssignedTo)
		if err == nil {
			assignedTo = &parsedID
		}
	}

	if err := h.service.UpdateTask(r.Context(), id, input.Title, input.Description, input.Status, input.Priority, assignedTo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Delete godoc
// @Summary Delete task
// @Description Delete a task by ID
// @Tags tasks
// @Param id path string true "Task ID"
// @Success 204 {string} string "No Content"
// @Router /api/tasks/{id} [delete]
func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteTask(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Stats godoc
// @Summary Get task statistics
// @Description Get various statistics about tasks
// @Tags tasks
// @Success 200 {object} models.TaskStats
// @Router /api/tasks/stats [get]
func (h *TaskHandler) Stats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// ExportCSV godoc
// @Summary Export tasks to CSV
// @Description Download all tasks in CSV format
// @Tags tasks
// @Produce text/csv
// @Success 200 {string} string "CSV data"
// @Router /api/tasks/export [get]
func (h *TaskHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.ListAllTasks(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=tasks.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{"ID", "Title", "Description", "Status", "Priority", "Created At"})
	for _, t := range tasks {
		writer.Write([]string{
			t.ID.String(),
			t.Title,
			t.Description,
			t.Status,
			t.Priority,
			t.CreatedAt.Format(time.RFC3339),
		})
	}
}
