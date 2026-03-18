package postgres

import (
	"context"

	"strconv"
 
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/user/devpulse/internal/models"
	"github.com/user/devpulse/internal/repository"
)

type taskRepo struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) repository.TaskRepository {
	return &taskRepo{pool: pool}
}

func (r *taskRepo) Create(ctx context.Context, t *models.Task) error {
	query := `INSERT INTO tasks (project_id, title, description, status, priority, assigned_to, due_date) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`
	return r.pool.QueryRow(ctx, query, t.ProjectID, t.Title, t.Description, t.Status, t.Priority, t.AssignedTo, t.DueDate).
		Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *taskRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	t := &models.Task{}
	query := `SELECT t.id, t.project_id, t.title, t.description, t.status, t.priority, t.assigned_to, u.full_name, t.due_date, t.created_at, t.updated_at 
              FROM tasks t 
              LEFT JOIN users u ON t.assigned_to = u.id 
              WHERE t.id = $1`
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.AssignedTo, &t.AssignedToName, &t.DueDate, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *taskRepo) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*models.Task, error) {
	query := `SELECT t.id, t.project_id, t.title, t.description, t.status, t.priority, t.assigned_to, u.full_name, t.due_date, t.created_at, t.updated_at 
              FROM tasks t 
              LEFT JOIN users u ON t.assigned_to = u.id 
              WHERE t.project_id = $1`
	rows, err := r.pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		t := &models.Task{}
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.AssignedTo, &t.AssignedToName, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *taskRepo) Update(ctx context.Context, t *models.Task) error {
	query := `UPDATE tasks SET title = $1, description = $2, status = $3, priority = $4, assigned_to = $5, due_date = $6, updated_at = CURRENT_TIMESTAMP WHERE id = $7`
	_, err := r.pool.Exec(ctx, query, t.Title, t.Description, t.Status, t.Priority, t.AssignedTo, t.DueDate, t.ID)
	return err
}

func (r *taskRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *taskRepo) GetStats(ctx context.Context) (map[string]int, error) {
	query := `SELECT status, COUNT(*) FROM tasks GROUP BY status`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		stats[status] = count
	}
	return stats, nil
}

func (r *taskRepo) ListAll(ctx context.Context) ([]*models.Task, error) {
	query := `SELECT t.id, t.project_id, t.title, t.description, t.status, t.priority, t.assigned_to, u.full_name, t.due_date, t.created_at, t.updated_at 
              FROM tasks t 
              LEFT JOIN users u ON t.assigned_to = u.id`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		t := &models.Task{}
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.AssignedTo, &t.AssignedToName, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
func (r *taskRepo) Search(ctx context.Context, query string) ([]*models.Task, error) {
	sqlQuery := `SELECT t.id, t.project_id, t.title, t.description, t.status, t.priority, t.assigned_to, u.full_name, t.due_date, t.created_at, t.updated_at 
                  FROM tasks t 
                  LEFT JOIN users u ON t.assigned_to = u.id 
                  WHERE t.title ILIKE $1 OR t.description ILIKE $1`
	rows, err := r.pool.Query(ctx, sqlQuery, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		t := &models.Task{}
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.AssignedTo, &t.AssignedToName, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
func (r *taskRepo) ListFiltered(ctx context.Context, projectID *uuid.UUID, priority string) ([]*models.Task, error) {
	query := `SELECT t.id, t.project_id, t.title, t.description, t.status, t.priority, t.assigned_to, u.full_name, t.due_date, t.created_at, t.updated_at 
              FROM tasks t 
              LEFT JOIN users u ON t.assigned_to = u.id 
              WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if projectID != nil {
		query += ` AND t.project_id = $` + strconv.Itoa(argIdx)
		args = append(args, *projectID)
		argIdx++
	}

	if priority != "" {
		query += ` AND t.priority = $` + strconv.Itoa(argIdx)
		args = append(args, priority)
		argIdx++
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		t := &models.Task{}
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.AssignedTo, &t.AssignedToName, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
