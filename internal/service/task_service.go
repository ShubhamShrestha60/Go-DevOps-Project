package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/user/devpulse/internal/models"
	"github.com/user/devpulse/internal/repository"
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, projectID uuid.UUID, title, description, priority string, assignedTo uuid.UUID) (*models.Task, error) {
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
	return task, nil
}

func (s *TaskService) GetTask(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) ListProjectTasks(ctx context.Context, projectID uuid.UUID) ([]*models.Task, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *TaskService) UpdateTask(ctx context.Context, id uuid.UUID, title, description, status, priority string, assignedTo uuid.UUID) error {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	t.Title = title
	t.Description = description
	t.Status = status
	t.Priority = priority
	t.AssignedTo = assignedTo
	return s.repo.Update(ctx, t)
}

func (s *TaskService) DeleteTask(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *TaskService) GetStats(ctx context.Context) (map[string]interface{}, error) {
	count, err := s.repo.GetStats(ctx)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"total_tasks": count,
	}, nil
}

func (s *TaskService) ListAllTasks(ctx context.Context) ([]*models.Task, error) {
	return s.repo.ListAll(ctx)
}
