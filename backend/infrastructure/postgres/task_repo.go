package postgres

import (
	"context"
	"database/sql"
	"taskmanager/domain"
	"taskmanager/repository"
	"time"
)

type taskRepo struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) repository.TaskRepository {
	return &taskRepo{db: db}
}

func (r *taskRepo) CreateTask(ctx context.Context, task *domain.Task) error {
	query := `INSERT INTO tasks (user_id, title, description, deadline, status, category_id, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return r.db.QueryRowContext(ctx, query, task.UserID, task.Title, task.Description, task.Deadline, task.Status, task.CategoryID, task.CreatedAt).Scan(&task.ID)
}

func (r *taskRepo) CreateCheckpoints(ctx context.Context, taskID int, checkpoints []domain.Checkpoint) error {
	query := `INSERT INTO checkpoints (task_id, title, "order", completed) VALUES ($1, $2, $3, $4) RETURNING id`

	// Start transaction or just iterate (iterate for simplicity)
	for i := range checkpoints {
		err := r.db.QueryRowContext(ctx, query, taskID, checkpoints[i].Title, checkpoints[i].Order, checkpoints[i].Completed).Scan(&checkpoints[i].ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *taskRepo) GetTasksByUserID(ctx context.Context, userID int) ([]domain.Task, error) {
	query := `
		SELECT t.id, t.user_id, t.title, t.description, t.deadline, t.status, t.category_id, t.created_at, 
		       c.id, c.name 
		FROM tasks t 
		LEFT JOIN categories c ON t.category_id = c.id 
		WHERE t.user_id = $1 ORDER BY t.created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var t domain.Task
		var catID sql.NullInt64
		var catName sql.NullString
		
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Deadline, &t.Status, &t.CategoryID, &t.CreatedAt, &catID, &catName); err != nil {
			return nil, err
		}

		if catID.Valid {
			t.Category = &domain.Category{
				ID:   int(catID.Int64),
				Name: catName.String,
			}
		}

		// Fetch checkpoints
		cpQuery := `SELECT id, task_id, title, "order", completed, completed_at FROM checkpoints WHERE task_id = $1 ORDER BY "order" ASC`
		cpRows, err := r.db.QueryContext(ctx, cpQuery, t.ID)
		if err == nil {
			var checkpts []domain.Checkpoint
			for cpRows.Next() {
				var cp domain.Checkpoint
				cpRows.Scan(&cp.ID, &cp.TaskID, &cp.Title, &cp.Order, &cp.Completed, &cp.CompletedAt)
				checkpts = append(checkpts, cp)
			}
			cpRows.Close()
			t.Checkpoints = checkpts
		}

		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *taskRepo) GetTaskByID(ctx context.Context, taskID int) (*domain.Task, error) {
	query := `
		SELECT t.id, t.user_id, t.title, t.description, t.deadline, t.status, t.category_id, t.created_at,
		       c.id, c.name
		FROM tasks t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.id = $1`
	t := &domain.Task{}
	var catID sql.NullInt64
	var catName sql.NullString

	err := r.db.QueryRowContext(ctx, query, taskID).Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Deadline, &t.Status, &t.CategoryID, &t.CreatedAt, &catID, &catName)
	if err != nil {
		return nil, err
	}

	if catID.Valid {
		t.Category = &domain.Category{
			ID:   int(catID.Int64),
			Name: catName.String,
		}
	}

	// Fetch checkpoints
	cpQuery := `SELECT id, task_id, title, "order", completed, completed_at FROM checkpoints WHERE task_id = $1 ORDER BY "order" ASC`
	cpRows, err := r.db.QueryContext(ctx, cpQuery, t.ID)
	if err == nil {
		var checkpts []domain.Checkpoint
		for cpRows.Next() {
			var cp domain.Checkpoint
			cpRows.Scan(&cp.ID, &cp.TaskID, &cp.Title, &cp.Order, &cp.Completed, &cp.CompletedAt)
			checkpts = append(checkpts, cp)
		}
		cpRows.Close()
		t.Checkpoints = checkpts
	}

	return t, nil
}

func (r *taskRepo) UpdateCheckpointStatus(ctx context.Context, checkpointID int, completed bool) error {
	query := `UPDATE checkpoints SET completed = $1, completed_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, completed, checkpointID)
	return err
}

func (r *taskRepo) UpdateTaskStatus(ctx context.Context, taskID int, status string) error {
	query := `UPDATE tasks SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, taskID)
	return err
}

func (r *taskRepo) UpdateTask(ctx context.Context, taskID int, title, description string, deadline time.Time, checkpoints []string, categoryID *int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	taskQuery := `UPDATE tasks SET title = $1, description = $2, deadline = $3, category_id = $4 WHERE id = $5`
	if _, err := tx.ExecContext(ctx, taskQuery, title, description, deadline, categoryID, taskID); err != nil {
		return err
	}

	for i, cpTitle := range checkpoints {
		cpQuery := `UPDATE checkpoints SET title = $1 WHERE task_id = $2 AND "order" = $3`
		if _, err := tx.ExecContext(ctx, cpQuery, cpTitle, taskID, i+1); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *taskRepo) DeleteTask(ctx context.Context, taskID int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, taskID)
	return err
}
