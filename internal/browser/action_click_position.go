package browser

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
	"github.com/go-rod/rod/lib/proto"
	"go.uber.org/zap"
)

// ClickAtPosition кликает по координатам (x, y) на странице
// Это самый надежный способ клика - работает как реальный пользователь
func (c *Controller) ClickAtPosition(ctx context.Context, x, y int) error {
	logger.Info(ctx, "🎯 Clicking at position", zap.Int("x", x), zap.Int("y", y))

	// Двигаем мышь плавно к позиции (как реальный пользователь)
	targetPoint := proto.Point{X: float64(x), Y: float64(y)}
	if err := c.page.Mouse.MoveLinear(targetPoint, 5); err != nil {
		logger.Warn(ctx, "⚠️ Mouse move failed, trying direct move", zap.Error(err))
		c.page.Mouse.MustMoveTo(float64(x), float64(y))
	}

	// Небольшая пауза как у реального пользователя
	time.Sleep(100 * time.Millisecond)

	// Кликаем
	if err := c.page.Mouse.Click(proto.InputMouseButtonLeft, 1); err != nil {
		logger.Error(ctx, "❌ Click at position failed", zap.Int("x", x), zap.Int("y", y), zap.Error(err))
		return fmt.Errorf("click at position failed: %w", err)
	}

	// Ждем стабилизации DOM после клика
	logger.Info(ctx, "⏳ Waiting for DOM to stabilize after click...")
	if err := c.page.Timeout(5 * time.Second).WaitStable(500 * time.Millisecond); err != nil {
		logger.Warn(ctx, "⚠️ WaitStable timeout (non-critical)", zap.Error(err))
	}
	time.Sleep(500 * time.Millisecond)

	logger.Info(ctx, "✅ Click at position completed", zap.Int("x", x), zap.Int("y", y))
	return nil
}
