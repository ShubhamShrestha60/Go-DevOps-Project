package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/user/devpulse/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListAll(ctx context.Context) ([]*models.User, error)
}

type ProjectRepository interface {
	Create(ctx context.Context, project *models.Project) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error)
	ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]*models.Project, error)
	Search(ctx context.Context, query string) ([]*models.Project, error)
	Update(ctx context.Context, project *models.Project) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetStats(ctx context.Context) (int, error)
}

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Task, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*models.Task, error)
	ListAll(ctx context.Context, ownerID uuid.UUID) ([]*models.Task, error)
	Search(ctx context.Context, ownerID uuid.UUID, query string) ([]*models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetStats(ctx context.Context) (map[string]int, error)
	ListFiltered(ctx context.Context, ownerID uuid.UUID, projectID *uuid.UUID, priority string) ([]*models.Task, error)
}

type ActivityLogRepository interface {
	Create(ctx context.Context, log *models.ActivityLog) error
	List(ctx context.Context, limit int) ([]*models.ActivityLog, error)
}

type CommentRepository interface {
	Create(ctx context.Context, comment *models.Comment) error
	ListByTask(ctx context.Context, taskID uuid.UUID) ([]*models.Comment, error)
}
