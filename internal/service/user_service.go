package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/user/devpulse/internal/models"
	"github.com/user/devpulse/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo         repository.UserRepository
	activityRepo repository.ActivityLogRepository
}

func NewUserService(repo repository.UserRepository, activityRepo repository.ActivityLogRepository) *UserService {
	return &UserService{repo: repo, activityRepo: activityRepo}
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context) ([]*models.User, error) {
	return s.repo.ListAll(ctx)
}

func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, fullName, email, password string) error {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if fullName != "" { u.FullName = fullName }
	if email != "" { u.Email = email }
	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.PasswordHash = string(hash)
	}

	if err := s.repo.Update(ctx, u); err != nil {
		return err
	}

	// Log Activity
	s.activityRepo.Create(ctx, &models.ActivityLog{
		UserID:     id,
		Action:     "update",
		EntityType: "user",
		EntityID:   id,
		Details:    "Updated profile settings",
	})

	return nil
}
