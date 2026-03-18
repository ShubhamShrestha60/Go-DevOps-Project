package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/user/devpulse/internal/models"
	"github.com/user/devpulse/internal/repository"
)

type userRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) repository.UserRepository {
	return &userRepo{pool: pool}
}

func (r *userRepo) Create(ctx context.Context, u *models.User) error {
	query := `INSERT INTO users (username, email, password_hash, full_name, role) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	return r.pool.QueryRow(ctx, query, u.Username, u.Email, u.PasswordHash, u.FullName, u.Role).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	u := &models.User{}
	query := `SELECT id, username, email, password_hash, full_name, role, created_at, updated_at FROM users WHERE id = $1`
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}
	query := `SELECT id, username, email, password_hash, full_name, role, created_at, updated_at FROM users WHERE email = $1`
	err := r.pool.QueryRow(ctx, query, email).
		Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepo) Update(ctx context.Context, u *models.User) error {
	query := `UPDATE users SET username = $1, email = $2, full_name = $3, role = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5`
	_, err := r.pool.Exec(ctx, query, u.Username, u.Email, u.FullName, u.Role, u.ID)
	return err
}

func (r *userRepo) ListAll(ctx context.Context) ([]*models.User, error) {
	query := `SELECT id, username, email, full_name, role, created_at, updated_at FROM users`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
func (r *userRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
