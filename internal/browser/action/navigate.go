package action

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
	"go.uber.org/zap"
)

// Navigate переходит по URL
func Navigate(ctx context.Context, p PageProvider, url string) error {
	logger.Info(ctx, "🌐 Navigating", zap.String("url", url))

	page := p.GetPage()
	timeout := p.GetTimeout()

	if err := page.Timeout(timeout).Navigate(url); err != nil {
		logger.Error(ctx, "❌ Navigation failed", zap.Error(err))
		return fmt.Errorf("navigation failed: %w", err)
	}

	if err := page.Timeout(timeout).WaitLoad(); err != nil {
		logger.Error(ctx, "❌ WaitLoad failed", zap.Error(err))
		return fmt.Errorf("page load timeout: %w", err)
	}

	if err := page.Timeout(5 * time.Second).WaitStable(500 * time.Millisecond); err != nil {
		logger.Warn(ctx, "⚠️ WaitStable timeout", zap.Error(err))
	}

	logger.Info(ctx, "✅ Navigation completed", zap.String("url", url))
	return nil
}
