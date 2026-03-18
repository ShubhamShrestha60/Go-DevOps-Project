package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/user/devpulse/internal/models"
	"github.com/user/devpulse/internal/repository"
)

type TaskService struct {
	repo         repository.TaskRepository
	projectRepo  repository.ProjectRepository
	activityRepo repository.ActivityLogRepository
}

func NewTaskService(repo repository.TaskRepository, projectRepo repository.ProjectRepository, activityRepo repository.ActivityLogRepository) *TaskService {
	return &TaskService{repo: repo, projectRepo: projectRepo, activityRepo: activityRepo}
}

func (s *TaskService) CreateTask(ctx context.Context, userID, projectID uuid.UUID, title, description, priority string, assignedTo *uuid.UUID) (*models.Task, error) {
	// Verify project ownership
	proj, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if proj.OwnerID != userID {
		return nil, models.ErrUnauthorized
	}

	task := &models.Task{
		ProjectID:   projectID,
		Title:       title,
		Description: description,
		Status:      "todo",
		Priority:    priority,
		AssignedTo:  assignedTo,
	}
	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}
	s.activityRepo.Create(ctx, &models.ActivityLog{
		UserID:     userID,
		Action:     "create",
		EntityType: "task",
		EntityID:   task.ID,
		Details:    "Created task: " + title,
	})
	return task, nil
}

func (s *TaskService) GetTask(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) ListProjectTasks(ctx context.Context, projectID uuid.UUID) ([]*models.Task, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *TaskService) UpdateTask(ctx context.Context, userID, id uuid.UUID, title, description, status, priority string, assignedTo *uuid.UUID) error {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Verify project ownership
	proj, err := s.projectRepo.GetByID(ctx, t.ProjectID)
	if err != nil {
		return err
	}
	if proj.OwnerID != userID {
		return models.ErrUnauthorized
	}

	if title != "" { t.Title = title }
	if description != "" { t.Description = description }
	if status != "" { t.Status = status }
	if priority != "" { t.Priority = priority }
	if assignedTo != nil { t.AssignedTo = assignedTo }

	err = s.repo.Update(ctx, t)
	if err == nil {
		s.activityRepo.Create(ctx, &models.ActivityLog{
			UserID:     userID,
			Action:     "update",
			EntityType: "task",
			EntityID:   id,
			Details:    "Updated task: " + t.Title + " to status: " + t.Status,
		})
	}
	return err
}

func (s *TaskService) DeleteTask(ctx context.Context, userID, id uuid.UUID) error {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Verify project ownership
	proj, err := s.projectRepo.GetByID(ctx, t.ProjectID)
	if err != nil {
		return err
	}
	if proj.OwnerID != userID {
		return models.ErrUnauthorized
	}

	err = s.repo.Delete(ctx, id)
	if err == nil {
		s.activityRepo.Create(ctx, &models.ActivityLog{
			UserID:     userID,
			Action:     "delete",
			EntityType: "task",
			EntityID:   id,
			Details:    "Deleted task: " + t.Title,
		})
	}
	return err
}

func (s *TaskService) GetStats(ctx context.Context) (map[string]interface{}, error) {
	statusStats, err := s.repo.GetStats(ctx)
	if err != nil {
		return nil, err
	}
	
	total := 0
	for _, count := range statusStats {
		total += count
	}

	return map[string]interface{}{
		"total_tasks": total,
		"statuses":    statusStats,
	}, nil
}

func (s *TaskService) SearchTasks(ctx context.Context, userID uuid.UUID, query string) ([]*models.Task, error) {
	return s.repo.Search(ctx, userID, query)
}

func (s *TaskService) ListAllTasks(ctx context.Context, userID uuid.UUID) ([]*models.Task, error) {
	return s.repo.ListAll(ctx, userID)
}

func (s *TaskService) ListFiltered(ctx context.Context, userID uuid.UUID, projectID *uuid.UUID, priority string) ([]*models.Task, error) {
	return s.repo.ListFiltered(ctx, userID, projectID, priority)
}
