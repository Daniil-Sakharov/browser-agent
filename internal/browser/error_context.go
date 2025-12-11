package browser

import (
	"context"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/browser/dom"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
)

// BuildErrorContext создаёт контекст ошибки (делегирует в dom пакет)
func (c *Controller) BuildErrorContext(ctx context.Context, failedSelector string, err error) *domain.ErrorContext {
	return dom.BuildErrorContext(ctx, c, failedSelector, err)
}

// FormatErrorContextMessage форматирует сообщение об ошибке
func FormatErrorContextMessage(ctx *domain.ErrorContext) string {
	return dom.FormatErrorContextMessage(ctx)
}
