package browser

import (
	"context"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/browser/action"
)

// Click кликает по элементу (делегирует в action пакет)
func (c *Controller) Click(ctx context.Context, selector string) error {
	return action.Click(ctx, c, selector)
}
