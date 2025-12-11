package browser

import (
	"context"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/browser/action"
)

// PressEnter нажимает Enter (делегирует в action пакет)
func (c *Controller) PressEnter(ctx context.Context) error {
	return action.PressEnter(ctx, c)
}
