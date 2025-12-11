package browser

import (
	"context"
	"fmt"
	"os"

	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
	"go.uber.org/zap"
)

// GetURL возвращает текущий URL
func (c *Controller) GetURL(ctx context.Context) string {
	info := c.page.MustInfo()
	return info.URL
}

// GetTitle возвращает заголовок страницы
func (c *Controller) GetTitle(ctx context.Context) string {
	info := c.page.MustInfo()
	return info.Title
}

// Screenshot делает скриншот страницы
func (c *Controller) Screenshot(ctx context.Context, path string) error {
	data, err := c.page.Screenshot(true, nil)
	if err != nil {
		return fmt.Errorf("screenshot failed: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to save screenshot: %w", err)
	}

	logger.Info(ctx, "📸 Screenshot saved", zap.String("path", path))
	return nil
}
