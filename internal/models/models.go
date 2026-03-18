package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Project struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     uuid.UUID `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	TaskCount   int       `json:"task_count"`
}

type Task struct {
	ID          uuid.UUID `json:"id"`
	ProjectID   uuid.UUID `json:"project_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`   // todo, in-progress, review, done
	Priority    string    `json:"priority"` // low, medium, high, urgent
	AssignedTo      *uuid.UUID `json:"assigned_to,omitempty"`
	AssignedToName  string     `json:"assigned_to_name,omitempty"`
	DueDate         *time.Time `json:"due_date,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	ProjectName     string     `json:"project_name,omitempty"`
}

type ActivityLog struct {
	ID         uuid.UUID   `json:"id"`
	UserID     uuid.UUID   `json:"user_id"`
	UserName   string      `json:"user_name,omitempty"`
	Action     string      `json:"action"`
	EntityType string      `json:"entity_type"`
	EntityID   uuid.UUID   `json:"entity_id"`
	Details    interface{} `json:"details"`
	CreatedAt  time.Time   `json:"created_at"`
}

type Comment struct {
	ID        uuid.UUID `json:"id"`
	TaskID    uuid.UUID `json:"task_id"`
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type TaskStats struct {
	TotalTasks int            `json:"total_tasks"`
	Statuses   map[string]int `json:"statuses"`
}

type ProjectStats struct {
	TotalProjects int `json:"total_projects"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
