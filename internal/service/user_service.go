package service

import (
	"context"

	"github.com/user/devpulse/internal/models"
	"github.com/user/devpulse/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) ListUsers(ctx context.Context) ([]*models.User, error) {
	return s.repo.ListAll(ctx)
}
