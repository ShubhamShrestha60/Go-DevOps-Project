package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/user/devpulse/internal/models"
	"github.com/user/devpulse/internal/repository"
)

type projectRepo struct {
	pool *pgxpool.Pool
}

func NewProjectRepository(pool *pgxpool.Pool) repository.ProjectRepository {
	return &projectRepo{pool: pool}
}

func (r *projectRepo) Create(ctx context.Context, p *models.Project) error {
	query := `INSERT INTO projects (name, description, owner_id) 
              VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	return r.pool.QueryRow(ctx, query, p.Name, p.Description, p.OwnerID).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *projectRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	p := &models.Project{}
	query := `SELECT p.id, p.name, COALESCE(p.description, ''), p.owner_id, p.created_at, p.updated_at, COUNT(t.id) as task_count 
              FROM projects p 
              LEFT JOIN tasks t ON p.id = t.project_id 
              WHERE p.id = $1 
              GROUP BY p.id`
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID, &p.CreatedAt, &p.UpdatedAt, &p.TaskCount)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *projectRepo) ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]*models.Project, error) {
	query := `SELECT p.id, p.name, COALESCE(p.description, ''), p.owner_id, p.created_at, p.updated_at, COUNT(t.id) as task_count 
              FROM projects p 
              LEFT JOIN tasks t ON p.id = t.project_id 
              WHERE p.owner_id = $1 
              GROUP BY p.id
              ORDER BY p.created_at DESC`
	rows, err := r.pool.Query(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		p := &models.Project{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID, &p.CreatedAt, &p.UpdatedAt, &p.TaskCount); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (r *projectRepo) Update(ctx context.Context, p *models.Project) error {
	query := `UPDATE projects SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`
	_, err := r.pool.Exec(ctx, query, p.Name, p.Description, p.ID)
	return err
}

func (r *projectRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *projectRepo) GetStats(ctx context.Context) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM projects").Scan(&count)
	return count, err
}

func (r *projectRepo) Search(ctx context.Context, query string) ([]*models.Project, error) {
	sqlQuery := `SELECT id, name, COALESCE(description, ''), owner_id, created_at, updated_at FROM projects 
                 WHERE name ILIKE $1 OR description ILIKE $1`
	rows, err := r.pool.Query(ctx, sqlQuery, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		p := &models.Project{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}
