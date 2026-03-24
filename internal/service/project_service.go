package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/user/devpulse/internal/models"
	"github.com/user/devpulse/internal/repository"
)

type ProjectService struct {
	repo         repository.ProjectRepository
	activityRepo repository.ActivityLogRepository
}

func NewProjectService(repo repository.ProjectRepository, activityRepo repository.ActivityLogRepository) *ProjectService {
	return &ProjectService{repo: repo, activityRepo: activityRepo}
}

func (s *ProjectService) CreateProject(ctx context.Context, name, description string, ownerID uuid.UUID) (*models.Project, error) {
	project := &models.Project{
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
	}
	if err := s.repo.Create(ctx, project); err != nil {
		return nil, err
	}
	_ = s.activityRepo.Create(ctx, &models.ActivityLog{
		UserID:     ownerID,
		Action:     "create",
		EntityType: "project",
		EntityID:   project.ID,
		Details:    "Created project: " + name,
	})
	return project, nil
}

func (s *ProjectService) GetProject(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProjectService) GetStats(ctx context.Context) (map[string]interface{}, error) {
	count, err := s.repo.GetStats(ctx)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"total_projects": count,
	}, nil
}

func (s *ProjectService) ListUserProjects(ctx context.Context, userID uuid.UUID) ([]*models.Project, error) {
	return s.repo.ListByOwner(ctx, userID)
}

func (s *ProjectService) UpdateProject(ctx context.Context, userID, id uuid.UUID, name, description string) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if p.OwnerID != userID {
		return models.ErrUnauthorized
	}

	p.Name = name
	p.Description = description
	err = s.repo.Update(ctx, p)
	if err == nil {
		_ = s.activityRepo.Create(ctx, &models.ActivityLog{
			UserID:     p.OwnerID,
			Action:     "update",
			EntityType: "project",
			EntityID:   id,
			Details:    "Updated project: " + name,
		})
	}
	return err
}

func (s *ProjectService) DeleteProject(ctx context.Context, userID, id uuid.UUID) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if p.OwnerID != userID {
		return models.ErrUnauthorized
	}

	err = s.repo.Delete(ctx, id)
	if err == nil {
		_ = s.activityRepo.Create(ctx, &models.ActivityLog{
			UserID:     p.OwnerID,
			Action:     "delete",
			EntityType: "project",
			EntityID:   id,
			Details:    "Deleted project: " + p.Name,
		})
	}
	return err
}

func (s *ProjectService) SearchProjects(ctx context.Context, query string) ([]*models.Project, error) {
	return s.repo.Search(ctx, query)
}
