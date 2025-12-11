package browser

import (
	"context"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/browser/action"
)

// Scroll прокручивает страницу (делегирует в action пакет)
func (c *Controller) Scroll(ctx context.Context, direction string, amount int) error {
	return action.Scroll(ctx, c, direction, amount)
}
