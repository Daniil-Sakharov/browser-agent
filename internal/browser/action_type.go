package browser

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
	"go.uber.org/zap"
)

// Type вводит текст в поле
func (c *Controller) Type(ctx context.Context, selector, text string) error {
	logger.Info(ctx, "⌨️  Typing text",
		zap.String("selector", selector),
		zap.Int("length", len(text)))

	elementTimeout := 10 * time.Second

	elem, err := c.page.Timeout(elementTimeout).Element(selector)
	if err != nil {
		logger.Error(ctx, "❌ Element not found", zap.String("selector", selector), zap.Error(err))
		return fmt.Errorf("element not found: %w", err)
	}

	// Очищаем поле перед вводом
	if err := elem.SelectAllText(); err != nil {
		logger.Warn(ctx, "⚠️ Failed to select text", zap.Error(err))
	}

	// Вводим текст
	if err := elem.Input(text); err != nil {
		return fmt.Errorf("input failed: %w", err)
	}

	// Ждем стабилизации DOM после ввода (для автокомплита, валидации)
	if err := c.page.Timeout(3 * time.Second).WaitStable(300 * time.Millisecond); err != nil {
		// Игнорируем таймаут - не критично
	}

	logger.Info(ctx, "✅ Typing completed", zap.String("selector", selector))
	return nil
}
