package browser

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
)

// WaitForElement ждет появления элемента
func (c *Controller) WaitForElement(ctx context.Context, selector string, timeout time.Duration) error {
	logger.Info(ctx, "⏳ Waiting for element",
		zap.String("selector", selector),
		zap.Duration("timeout", timeout))

	if timeout == 0 {
		timeout = c.timeout
	}

	elem, err := c.page.Timeout(timeout).Element(selector)
	if err != nil {
		return fmt.Errorf("element wait timeout: %w", err)
	}

	// Ждем видимости
	if err := elem.Timeout(timeout).WaitVisible(); err != nil {
		return fmt.Errorf("element not visible: %w", err)
	}

	logger.Info(ctx, "✅ Element found", zap.String("selector", selector))
	return nil
}

// WaitForNavigation ждет завершения навигации
func (c *Controller) WaitForNavigation(ctx context.Context) error {
	logger.Info(ctx, "⏳ Waiting for navigation")

	if err := c.page.Timeout(c.timeout).WaitLoad(); err != nil {
		return fmt.Errorf("navigation timeout: %w", err)
	}

	logger.Info(ctx, "✅ Navigation completed")
	return nil
}

// Sleep простая пауза
func (c *Controller) Sleep(ctx context.Context, duration time.Duration) {
	logger.Info(ctx, "💤 Sleeping", zap.Duration("duration", duration))
	time.Sleep(duration)
}
