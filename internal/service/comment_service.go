package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/user/devpulse/internal/models"
	"github.com/user/devpulse/internal/repository"
)

type CommentService struct {
	repo         repository.CommentRepository
	activityRepo repository.ActivityLogRepository
}

func NewCommentService(repo repository.CommentRepository, activityRepo repository.ActivityLogRepository) *CommentService {
	return &CommentService{repo: repo, activityRepo: activityRepo}
}

func (s *CommentService) AddComment(ctx context.Context, taskID, userID uuid.UUID, content string) (*models.Comment, error) {
	comment := &models.Comment{
		TaskID:  taskID,
		UserID:  userID,
		Content: content,
	}
	if err := s.repo.Create(ctx, comment); err != nil {
		return nil, err
	}

	// Log Activity
	s.activityRepo.Create(ctx, &models.ActivityLog{
		UserID:     userID,
		Action:     "create",
		EntityType: "comment",
		EntityID:   comment.ID,
		Details:    "Added a comment to task",
	})

	return comment, nil
}

func (s *CommentService) GetTaskComments(ctx context.Context, taskID uuid.UUID) ([]*models.Comment, error) {
	return s.repo.ListByTask(ctx, taskID)
}
