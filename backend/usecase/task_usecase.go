package usecase

import (
	"context"
	"errors"
	"time"

	"taskmanager/domain"
	"taskmanager/repository"
)

type TaskUseCase interface {
	CreateTask(ctx context.Context, userID int, title, description string, deadline time.Time, checkpoints []string, categoryID *int) (*domain.Task, error)
	GetTasks(ctx context.Context, userID int) ([]domain.Task, error)
	CompleteCheckpoint(ctx context.Context, userID, taskID, checkpointID int, completed bool) error
	UpdateTask(ctx context.Context, userID, taskID int, title, description string, deadline time.Time, checkpoints []string, categoryID *int) error
	DeleteTask(ctx context.Context, userID, taskID int) error
}

type taskUseCase struct {
	taskRepo repository.TaskRepository
}

func NewTaskUseCase(taskRepo repository.TaskRepository) TaskUseCase {
	return &taskUseCase{
		taskRepo: taskRepo,
	}
}

func (u *taskUseCase) CreateTask(ctx context.Context, userID int, title, description string, deadline time.Time, cpTitles []string, categoryID *int) (*domain.Task, error) {
	if len(cpTitles) != 4 {
		return nil, errors.New("exactly 4 checkpoints are required")
	}

	task := &domain.Task{
		UserID:      userID,
		Title:       title,
		Description: description,
		Deadline:    deadline,
		Status:      "Not Started",
		CategoryID:  categoryID,
		CreatedAt:   time.Now(),
	}

	if err := u.taskRepo.CreateTask(ctx, task); err != nil {
		return nil, err
	}

	checkpoints := make([]domain.Checkpoint, 4)
	for i, cpTitle := range cpTitles {
		checkpoints[i] = domain.Checkpoint{
			TaskID: task.ID,
			Title:  cpTitle,
			Order:  i + 1,
		}
	}

	if err := u.taskRepo.CreateCheckpoints(ctx, task.ID, checkpoints); err != nil {
		return nil, err
	}

	task.Checkpoints = checkpoints
	return task, nil
}

func (u *taskUseCase) GetTasks(ctx context.Context, userID int) ([]domain.Task, error) {
	return u.taskRepo.GetTasksByUserID(ctx, userID)
}

func (u *taskUseCase) CompleteCheckpoint(ctx context.Context, userID, taskID, checkpointID int, completed bool) error {
	task, err := u.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil || task == nil {
		return errors.New("task not found")
	}

	if task.UserID != userID {
		return errors.New("unauthorized")
	}

	if task.Status == "Completed" && !completed {
		return errors.New("cannot untick checkpoints of a fully completed task")
	}

	// Find the target checkpoint to enforce sequential ticking
	var targetCp *domain.Checkpoint
	for _, cp := range task.Checkpoints {
		if cp.ID == checkpointID {
			targetCp = &cp
			break
		}
	}

	if targetCp == nil {
		return errors.New("checkpoint not found in task")
	}

	// Enforce sequential ticking/unticking
	if completed {
		// If we're completing this checkpoint, verify all prior ones are completed
		for _, cp := range task.Checkpoints {
			if cp.Order < targetCp.Order && !cp.Completed {
				return errors.New("cannot check this point before completing previous points")
			}
		}
	} else {
		// If we're unticking this checkpoint, verify no subsequent ones are completed
		for _, cp := range task.Checkpoints {
			if cp.Order > targetCp.Order && cp.Completed {
				return errors.New("cannot untick this point while subsequent points are completed")
			}
		}
	}

	// Update checkpoint
	if err := u.taskRepo.UpdateCheckpointStatus(ctx, checkpointID, completed); err != nil {
		return err
	}

	// Re-fetch to calculate new status
	updatedTask, err := u.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		return err
	}

	completedCount := 0
	for _, cp := range updatedTask.Checkpoints {
		if cp.Completed {
			completedCount++
		}
	}

	newStatus := "Not Started"
	switch completedCount {
	case 1, 2:
		newStatus = "In Progress"
	case 3:
		newStatus = "Almost Done"
	case 4:
		newStatus = "Completed"
	}

	if newStatus != task.Status {
		return u.taskRepo.UpdateTaskStatus(ctx, taskID, newStatus)
	}

	return nil
}

func (u *taskUseCase) DeleteTask(ctx context.Context, userID, taskID int) error {
	task, err := u.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil || task == nil {
		return errors.New("task not found")
	}

	if task.UserID != userID {
		return errors.New("unauthorized")
	}

	return u.taskRepo.DeleteTask(ctx, taskID)
}

func (u *taskUseCase) UpdateTask(ctx context.Context, userID, taskID int, title, description string, deadline time.Time, checkpoints []string, categoryID *int) error {
	if len(checkpoints) != 4 {
		return errors.New("exactly 4 checkpoints are required")
	}

	task, err := u.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil || task == nil {
		return errors.New("task not found")
	}

	if task.UserID != userID {
		return errors.New("unauthorized")
	}

	return u.taskRepo.UpdateTask(ctx, taskID, title, description, deadline, checkpoints, categoryID)
}
