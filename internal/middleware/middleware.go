package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/user/devpulse/internal/service"
	"go.uber.org/zap"
)

type Middleware struct {
	logger      *zap.Logger
	authService *service.AuthService
}

func New(logger *zap.Logger, authService *service.AuthService) *Middleware {
	return &Middleware{
		logger:      logger,
		authService: authService,
	}
}

// Logging middleware logs the incoming HTTP requests
func (m *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		m.logger.Info("request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Duration("duration", time.Since(start)),
		)
	})
}

// Auth middleware validates the JWT token in Authorization header or Cookie
func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// Check Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// If no header, check cookie
		if tokenString == "" {
			cookie, err := r.Cookie("token")
			if err == nil {
				tokenString = cookie.Value
			}
		}

		if tokenString == "" {
			m.handleUnauthorized(w, r)
			return
		}

		userID, role, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			m.handleUnauthorized(w, r)
			return
		}

		// Store user info in context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		ctx = context.WithValue(ctx, "user_role", role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) handleUnauthorized(w http.ResponseWriter, r *http.Request) {
	// If it's an API request, return JSON
	if strings.HasPrefix(r.URL.Path, "/api/") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = w.Write([]byte(`{"error": "unauthorized"}`))
		return
	}
	// Otherwise redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func GetUserID(ctx context.Context) uuid.UUID {
	id, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return id
}
