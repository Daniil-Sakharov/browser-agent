package action

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
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

	// WaitLoad - не критичная ошибка если таймаут
	if err := page.Timeout(timeout).WaitLoad(); err != nil {
		logger.Warn(ctx, "⚠️ WaitLoad timeout, continuing...", zap.Error(err))
		// Не возвращаем ошибку - страница может работать
	}

	// WaitStable опционально
	if err := page.Timeout(3 * time.Second).WaitStable(300 * time.Millisecond); err != nil {
		logger.Debug(ctx, "WaitStable timeout", zap.Error(err))
	}

	logger.Info(ctx, "✅ Navigation completed", zap.String("url", url))
	return nil
}
