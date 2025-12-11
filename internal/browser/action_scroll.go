package browser

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
	"go.uber.org/zap"
)

// Scroll прокручивает страницу
func (c *Controller) Scroll(ctx context.Context, direction string, amount int) error {
	logger.Info(ctx, "📜 Scrolling",
		zap.String("direction", direction),
		zap.Int("amount", amount))

	if amount == 0 {
		amount = 500
	}

	if direction == "down" {
		amount = -amount
	}

	err := c.page.Mouse.Scroll(0, float64(amount), 10)
	if err != nil {
		return fmt.Errorf("scroll failed: %w", err)
	}

	// Ждем стабилизации DOM (lazy loading)
	if err := c.page.Timeout(3 * time.Second).WaitStable(300 * time.Millisecond); err != nil {
		// Игнорируем таймаут
	}

	return nil
}
