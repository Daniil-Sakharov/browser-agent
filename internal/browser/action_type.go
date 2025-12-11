package browser

import (
	"context"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/browser/action"
)

// Type вводит текст в поле (делегирует в action пакет)
func (c *Controller) Type(ctx context.Context, selector, text string) error {
	return action.Type(ctx, c, selector, text)
}
