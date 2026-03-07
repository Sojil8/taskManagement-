package repository

import (
	"context"
	"time"

	"taskmanager/domain"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *domain.Task) error
	CreateCheckpoints(ctx context.Context, taskID int, checkpoints []domain.Checkpoint) error
	GetTasksByUserID(ctx context.Context, userID int) ([]domain.Task, error)
	GetTaskByID(ctx context.Context, taskID int) (*domain.Task, error)
	UpdateCheckpointStatus(ctx context.Context, checkpointID int, completed bool) error
	UpdateTask(ctx context.Context, taskID int, title, description string, deadline time.Time, checkpoints []string, categoryID *int) error
	UpdateTaskStatus(ctx context.Context, taskID int, status string) error
	DeleteTask(ctx context.Context, taskID int) error
}
