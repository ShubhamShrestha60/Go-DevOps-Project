package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/devpulse/internal/config"
	"github.com/user/devpulse/internal/database"
	"github.com/user/devpulse/internal/handler"
	"github.com/user/devpulse/internal/middleware"
	"github.com/user/devpulse/internal/repository/postgres"
	"github.com/user/devpulse/internal/router"
	"github.com/user/devpulse/internal/service"
	"go.uber.org/zap"
)

func main() {
	// 1. Load Config
	cfg := config.Load()

	// 2. Initialize Logger
	logger, _ := zap.NewProduction()
	if cfg.Env == "development" {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()

	// 3. Initialize Database
	db, err := database.New(cfg)
	if err != nil {
		logger.Fatal("could not initialize database", zap.Error(err))
	}
	defer db.Close()

	// 4. Initialize Repositories
	userRepo := postgres.NewUserRepository(db.Pool)
	projectRepo := postgres.NewProjectRepository(db.Pool)
	taskRepo := postgres.NewTaskRepository(db.Pool)

	// 5. Initialize Services
	authService := service.NewAuthService(userRepo, cfg.Auth.JWTSecret, cfg.Auth.JWTExpiryH, cfg.Auth.AdminPassword)
	projectService := service.NewProjectService(projectRepo)
	taskService := service.NewTaskService(taskRepo)

	// 6. Initialize Handlers
	authHandler := handler.NewAuthHandler(authService)
	projectHandler := handler.NewProjectHandler(projectService)
	taskHandler := handler.NewTaskHandler(taskService)
	healthHandler := handler.NewHealthHandler(db)
	dashboardHandler := handler.NewDashboardHandler()

	// 7. Initialize Middleware
	mw := middleware.New(logger, authService)

	// 8. Initialize Router
	r := router.New(mw, authHandler, projectHandler, taskHandler, healthHandler, dashboardHandler)

	// 9. Start Server with Graceful Shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		logger.Info("Starting server", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
