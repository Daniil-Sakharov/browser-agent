package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Ошибки Task
var (
	ErrTaskAlreadyStarted   = errors.New("task already started")
	ErrTaskAlreadyCompleted = errors.New("task already completed")
	ErrTaskNotRunning       = errors.New("task not running")
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

// Start запускает задачу
func (t *Task) Start() error {
	if t.Status != TaskStatusPending {
		return ErrTaskAlreadyStarted
	}
	t.Status = TaskStatusRunning
	return nil
}

// Complete завершает задачу успешно
func (t *Task) Complete(result string) error {
	if t.Status != TaskStatusRunning {
		return ErrTaskNotRunning
	}
	t.Status = TaskStatusCompleted
	t.Result = result
	now := time.Now()
	t.CompletedAt = &now
	return nil
}

// Fail завершает задачу с ошибкой
func (t *Task) Fail(err error) error {
	if t.Status == TaskStatusCompleted {
		return ErrTaskAlreadyCompleted
	}
	t.Status = TaskStatusFailed
	t.Error = err
	now := time.Now()
	t.CompletedAt = &now
	return nil
}

// IsRunning проверяет запущена ли задача
func (t *Task) IsRunning() bool {
	return t.Status == TaskStatusRunning
}

// IsCompleted проверяет завершена ли задача
func (t *Task) IsCompleted() bool {
	return t.Status == TaskStatusCompleted || t.Status == TaskStatusFailed
}

// TaskStatus статус выполнения задачи
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)
