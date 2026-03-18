package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/user/devpulse/internal/models"
	"github.com/user/devpulse/internal/repository"
)

type ActivityService struct {
	repo repository.ActivityLogRepository
}

func NewActivityService(repo repository.ActivityLogRepository) *ActivityService {
	return &ActivityService{repo: repo}
}

func (s *ActivityService) Log(ctx context.Context, userID uuid.UUID, action, entityType string, entityID uuid.UUID, details interface{}) error {
	log := &models.ActivityLog{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Details:    details,
	}
	return s.repo.Create(ctx, log)
}

func (s *ActivityService) GetRecent(ctx context.Context, limit int) ([]*models.ActivityLog, error) {
	return s.repo.List(ctx, limit)
}
