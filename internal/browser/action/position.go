package action

import (
	"context"
	"fmt"
	"time"

	"github.com/go-rod/rod/lib/proto"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
	"go.uber.org/zap"
)

// ClickAtPosition кликает по координатам (x, y)
func ClickAtPosition(ctx context.Context, p PageProvider, x, y int) error {
	logger.Info(ctx, "🎯 Clicking at position", zap.Int("x", x), zap.Int("y", y))

	page := p.GetPage()
	target := proto.Point{X: float64(x), Y: float64(y)}

	if err := page.Mouse.MoveLinear(target, 5); err != nil {
		page.Mouse.MustMoveTo(float64(x), float64(y))
	}

	time.Sleep(100 * time.Millisecond)

	if err := page.Mouse.Click(proto.InputMouseButtonLeft, 1); err != nil {
		logger.Error(ctx, "❌ Click failed", zap.Int("x", x), zap.Int("y", y), zap.Error(err))
		return fmt.Errorf("click at position failed: %w", err)
	}

	p.WaitStable(5 * time.Second)
	time.Sleep(500 * time.Millisecond)

	logger.Info(ctx, "✅ Click at position completed", zap.Int("x", x), zap.Int("y", y))
	return nil
}
