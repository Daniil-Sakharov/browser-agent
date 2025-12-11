package browser

import (
	"context"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/browser/action"
)

// ClickAtPosition кликает по координатам (делегирует в action пакет)
func (c *Controller) ClickAtPosition(ctx context.Context, x, y int) error {
	return action.ClickAtPosition(ctx, c, x, y)
}
