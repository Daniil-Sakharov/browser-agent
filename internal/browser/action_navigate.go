package browser

import (
	"context"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/browser/action"
)

// Navigate переходит по URL (делегирует в action пакет)
func (c *Controller) Navigate(ctx context.Context, url string) error {
	return action.Navigate(ctx, c, url)
}
