package browser

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
	"go.uber.org/zap"
)

// Navigate переходит по URL
func (c *Controller) Navigate(ctx context.Context, url string) error {
	logger.Info(ctx, "🌐 Navigating to URL", zap.String("url", url))

	if c.page == nil {
		return fmt.Errorf("page is nil")
	}

	// Таймаут на навигацию - 30 секунд
	err := c.page.Timeout(c.timeout).Navigate(url)
	if err != nil {
		logger.Error(ctx, "❌ Navigation failed", zap.String("url", url), zap.Error(err))
		return fmt.Errorf("navigation failed: %w", err)
	}

	logger.Info(ctx, "📄 Page.Navigate completed, waiting for load...")

	// Ждем загрузки страницы
	err = c.page.Timeout(c.timeout).WaitLoad()
	if err != nil {
		logger.Error(ctx, "❌ WaitLoad failed", zap.String("url", url), zap.Error(err))
		return fmt.Errorf("page load timeout: %w", err)
	}

	// Ждем стабилизации DOM (для SPA/React)
	logger.Info(ctx, "⏳ Waiting for page to stabilize...")
	if err := c.page.Timeout(5 * time.Second).WaitStable(500 * time.Millisecond); err != nil {
		logger.Warn(ctx, "⚠️ WaitStable timeout (non-critical)", zap.Error(err))
	}

	logger.Info(ctx, "✅ Navigation completed", zap.String("url", url))
	return nil
}
