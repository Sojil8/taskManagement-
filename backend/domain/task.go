package domain

import "time"

type Task struct {
	ID          int          `json:"id"`
	UserID      int          `json:"user_id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Deadline    time.Time    `json:"deadline"`
	Status      string       `json:"status"` // e.g., "Not Started", "In Progress", "Almost Done", "Completed"
	CreatedAt   time.Time    `json:"created_at"`
	Checkpoints []Checkpoint `json:"checkpoints,omitempty"`
	CategoryID  *int         `json:"category_id,omitempty"`
	Category    *Category    `json:"category,omitempty"`
}

type Checkpoint struct {
	ID          int        `json:"id"`
	TaskID      int        `json:"task_id"`
	Title       string     `json:"title"`
	Order       int        `json:"order"` // 1, 2, 3, 4
	Completed   bool       `json:"completed"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}
