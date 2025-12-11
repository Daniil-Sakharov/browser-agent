package agent

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
	"go.uber.org/zap"
)

// Agent основной orchestrator
type Agent struct {
	browser             BrowserController
	ai                  AIClient
	security            SecurityChecker
	domSubAgent         DOMSubAgent
	maxSteps            int
	interactive         bool
	screenshots         bool
	currentTask         *domain.Task
	stepCount           int
	consecutiveFailures int
	lastFailedAction    string
	progressCallback    ProgressCallback
}

// New создает новый Agent
func New(ctx context.Context, browser BrowserController, ai AIClient, security SecurityChecker, domSubAgent DOMSubAgent, maxSteps int, interactive, screenshots bool) (*Agent, error) {
	logger.Info(ctx, "✅ Agent initialized",
		zap.Int("max_steps", maxSteps),
		zap.Bool("interactive", interactive))

	return &Agent{
		browser: browser, ai: ai, security: security, domSubAgent: domSubAgent,
		maxSteps: maxSteps, interactive: interactive, screenshots: screenshots,
	}, nil
}

// SetProgressCallback устанавливает callback для вывода прогресса
func (a *Agent) SetProgressCallback(cb ProgressCallback) { a.progressCallback = cb }

// emitProgress отправляет событие прогресса
func (a *Agent) emitProgress(event ProgressEvent) {
	if a.progressCallback != nil {
		a.progressCallback(event)
	}
}

// Execute выполняет задачу (без лимита шагов)
func (a *Agent) Execute(ctx context.Context, task *domain.Task) error {
	logger.Info(ctx, "🚀 Starting task", zap.String("task_id", task.ID))

	a.currentTask = task
	a.stepCount = 0
	a.ai.NewConversation()

	if err := task.Start(); err != nil {
		return fmt.Errorf("start task: %w", err)
	}

	for {
		a.stepCount++
		logger.Info(ctx, "📍 Step", zap.Int("step", a.stepCount))

		complete, err := a.executeStep(ctx)
		if err != nil {
			_ = task.Fail(err)
			return fmt.Errorf("step %d: %w", a.stepCount, err)
		}

		if complete {
			_ = task.Complete(task.Result)
			logger.Info(ctx, "✅ Task completed", zap.Int("steps", a.stepCount))
			return nil
		}
	}
}

// Close закрывает агента
func (a *Agent) Close(ctx context.Context) error {
	logger.Info(ctx, "🚫 Closing Agent")
	return nil
}
