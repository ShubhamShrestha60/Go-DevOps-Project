package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/user/devpulse/internal/models"
	"github.com/user/devpulse/internal/repository"
)

type ProjectService struct {
	repo repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
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
	return project, nil
}

func (s *ProjectService) GetProject(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProjectService) ListUserProjects(ctx context.Context, userID uuid.UUID) ([]*models.Project, error) {
	return s.repo.ListByOwner(ctx, userID)
}

func (s *ProjectService) UpdateProject(ctx context.Context, id uuid.UUID, name, description string) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	p.Name = name
	p.Description = description
	return s.repo.Update(ctx, p)
}

func (s *ProjectService) DeleteProject(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
