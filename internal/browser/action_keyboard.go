package browser

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
	"github.com/go-rod/rod/lib/input"
)

// PressEnter нажимает Enter
func (c *Controller) PressEnter(ctx context.Context) error {
	logger.Info(ctx, "⏎ Pressing Enter")

	if err := c.page.Keyboard.Press(input.Enter); err != nil {
		return fmt.Errorf("press enter failed: %w", err)
	}

	// Ждем стабилизации DOM (отправка формы может вызвать изменения)
	logger.Info(ctx, "⏳ Waiting for DOM to stabilize after Enter...")
	if err := c.page.Timeout(5 * time.Second).WaitStable(500 * time.Millisecond); err != nil {
		logger.Warn(ctx, "⚠️ WaitStable timeout (non-critical)")
	}

	return nil
}
