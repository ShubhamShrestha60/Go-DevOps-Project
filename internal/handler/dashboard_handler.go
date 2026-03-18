package handler

import (
	"net/http"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

func (h *DashboardHandler) Index(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "dashboard", nil)
}

func (h *DashboardHandler) Projects(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "projects", nil)
}

func (h *DashboardHandler) Tasks(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "tasks", nil)
}

func (h *DashboardHandler) Profile(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "profile", nil)
}

func (h *DashboardHandler) Register(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/pages/register.html")
}
