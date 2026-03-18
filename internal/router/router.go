package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/user/devpulse/internal/handler"
	"github.com/user/devpulse/internal/middleware"
)

func New(
	mw *middleware.Middleware,
	authHandler *handler.AuthHandler,
	projectHandler *handler.ProjectHandler,
	taskHandler *handler.TaskHandler,
	userHandler *handler.UserHandler,
	activityHandler *handler.ActivityHandler,
	healthHandler *handler.HealthHandler,
	dashboardHandler *handler.DashboardHandler,
) *chi.Mux {
	r := chi.NewRouter()

	// Base middleware
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(mw.Logging)
	r.Use(chiMiddleware.Recoverer)

	// Static files
	fs := http.FileServer(http.Dir("web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Health and Metrics (Public)
	r.Get("/healthz", healthHandler.Healthz)
	r.Get("/readyz", healthHandler.Readyz)
	r.Handle("/metrics", promhttp.Handler())

	// Auth (Public)
	r.Post("/api/auth/register", authHandler.Register)
	r.Post("/api/auth/login", authHandler.Login)
	r.Post("/api/auth/logout", authHandler.Logout)

	// Protected API
	r.Group(func(r chi.Router) {
		r.Use(mw.Auth)

		r.Route("/api/projects", func(r chi.Router) {
			r.Get("/stats", projectHandler.Stats)
			r.Post("/", projectHandler.Create)
			r.Get("/", projectHandler.List)
			r.Get("/{id}", projectHandler.Get)
			r.Put("/{id}", projectHandler.Update)
			r.Delete("/{id}", projectHandler.Delete)
		})

		r.Route("/api/tasks", func(r chi.Router) {
			r.Get("/stats", taskHandler.Stats)
			r.Post("/", taskHandler.Create)
			r.Get("/", taskHandler.List)
			r.Get("/{id}", taskHandler.Get)
			r.Put("/{id}", taskHandler.Update)
			r.Delete("/{id}", taskHandler.Delete)
		})

		r.Get("/api/users", userHandler.List)
		r.Get("/api/activities", activityHandler.List)
	})

	// Frontend Routes (Placeholder for template rendering)
	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/pages/login.html")
	})

	r.Group(func(r chi.Router) {
		r.Use(mw.Auth)
		r.Get("/", dashboardHandler.Index)
		r.Get("/projects", dashboardHandler.Projects)
		r.Get("/tasks", dashboardHandler.Tasks)
	})

	return r
}
