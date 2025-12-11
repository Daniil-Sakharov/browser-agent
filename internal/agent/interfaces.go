package agent

import (
	"context"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
)

// BrowserController интерфейс браузера
type BrowserController interface {
	GetPageContext(ctx context.Context) (*domain.PageContext, error)
	ExecuteAction(ctx context.Context, action domain.Action) (*domain.ActionResult, error)
	GetHTML(ctx context.Context) (string, error)
	FindElementsLive(ctx context.Context, query string) (string, error)
	Close(ctx context.Context) error
}

// AIClient интерфейс AI клиента
type AIClient interface {
	NewConversation()
	AddUserMessage(task string, pageContext *domain.PageContext) error
	AddToolResult(toolUseID, result string, isError bool)
	AddToolResultWithImage(toolUseID, result, imageB64 string, isError bool)
	DecideNextAction(ctx context.Context) (*domain.Decision, error)
	Close(ctx context.Context) error
}

// DOMSubAgent интерфейс для DOM sub-agent
type DOMSubAgent interface {
	Analyze(ctx context.Context, html, liveElements, question string) (string, error)
	AnalyzeError(ctx context.Context, html, liveElements, failedAction, errorMsg string) (string, error)
}

// SecurityChecker интерфейс проверки безопасности
type SecurityChecker interface {
	CheckAction(ctx context.Context, action domain.Action, pageContext *domain.PageContext) error
	Close(ctx context.Context) error
}

// ProgressCallback функция для вывода прогресса
type ProgressCallback func(event ProgressEvent)

// ProgressEvent событие прогресса
type ProgressEvent struct {
	Type      string
	Step      int
	MaxSteps  int
	Reasoning string
	Tool      string
	Params    map[string]string
	Result    string
	Success   bool
}
