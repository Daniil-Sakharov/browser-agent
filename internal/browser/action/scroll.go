package action

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
	"go.uber.org/zap"
)

// Scroll прокручивает страницу
func Scroll(ctx context.Context, p PageProvider, direction string, amount int) error {
	logger.Info(ctx, "📜 Scrolling", zap.String("direction", direction), zap.Int("amount", amount))

	if amount == 0 {
		amount = 500
	}
	if direction == "down" {
		amount = -amount
	}

	page := p.GetPage()
	if err := page.Mouse.Scroll(0, float64(amount), 10); err != nil {
		return fmt.Errorf("scroll failed: %w", err)
	}

	page.Timeout(3 * time.Second).WaitStable(300 * time.Millisecond)
	return nil
}
