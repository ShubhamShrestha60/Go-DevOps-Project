package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/user/devpulse/internal/models"
	"github.com/user/devpulse/internal/repository"
)

type activityRepo struct {
	pool *pgxpool.Pool
}

func NewActivityRepository(pool *pgxpool.Pool) repository.ActivityLogRepository {
	return &activityRepo{pool: pool}
}

func (r *activityRepo) Create(ctx context.Context, l *models.ActivityLog) error {
	query := `INSERT INTO activity_log (user_id, action, entity_type, entity_id, details) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query, l.UserID, l.Action, l.EntityType, l.EntityID, l.Details).
		Scan(&l.ID, &l.CreatedAt)
}

func (r *activityRepo) List(ctx context.Context, limit int) ([]*models.ActivityLog, error) {
	query := `SELECT l.id, l.user_id, COALESCE(u.full_name, u.username, 'System'), l.action, l.entity_type, l.entity_id, l.details, l.created_at 
              FROM activity_log l 
              LEFT JOIN users u ON l.user_id = u.id 
              ORDER BY l.created_at DESC LIMIT $1`
	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.ActivityLog
	for rows.Next() {
		l := &models.ActivityLog{}
		if err := rows.Scan(&l.ID, &l.UserID, &l.UserName, &l.Action, &l.EntityType, &l.EntityID, &l.Details, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}
