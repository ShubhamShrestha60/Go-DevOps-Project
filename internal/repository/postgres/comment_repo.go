package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/user/devpulse/internal/models"
	"github.com/user/devpulse/internal/repository"
)

type commentRepo struct {
	pool *pgxpool.Pool
}

func NewCommentRepository(pool *pgxpool.Pool) repository.CommentRepository {
	return &commentRepo{pool: pool}
}

func (r *commentRepo) Create(ctx context.Context, c *models.Comment) error {
	query := `INSERT INTO comments (task_id, user_id, content) 
              VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query, c.TaskID, c.UserID, c.Content).
		Scan(&c.ID, &c.CreatedAt)
}

func (r *commentRepo) ListByTask(ctx context.Context, taskID uuid.UUID) ([]*models.Comment, error) {
	query := `SELECT c.id, c.task_id, c.user_id, u.username, c.content, c.created_at 
              FROM comments c 
              JOIN users u ON c.user_id = u.id 
              WHERE c.task_id = $1 
              ORDER BY c.created_at ASC`
	rows, err := r.pool.Query(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		c := &models.Comment{}
		if err := rows.Scan(&c.ID, &c.TaskID, &c.UserID, &c.UserName, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}
