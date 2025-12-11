package domain

import (
	"time"

	"github.com/google/uuid"
)

// Task представляет задачу пользователя для агента
type Task struct {
	ID          string
	Description string
	Status      TaskStatus
	CreatedAt   time.Time
	CompletedAt *time.Time
	Result      string
	Error       error
}

// NewTask создает новую задачу
func NewTask(description string) *Task {
	return &Task{
		ID:          uuid.New().String(),
		Description: description,
		Status:      TaskStatusPending,
		CreatedAt:   time.Now(),
	}
}

// TaskStatus статус выполнения задачи
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)
