package browser

import (
	"context"
	"time"

	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
	"go.uber.org/zap"
)

// SwitchToNewTab переключается на новую вкладку если она появилась
func (c *Controller) SwitchToNewTab(ctx context.Context) bool {
	pages, err := c.browser.Pages()
	if err != nil {
		return false
	}

	if len(pages) <= 1 {
		return false
	}

	// Находим последнюю (новую) вкладку
	newPage := pages[len(pages)-1]
	if newPage == c.page {
		return false
	}

	// Переключаемся на новую вкладку
	c.page = newPage
	c.page.Timeout(c.timeout).WaitLoad()
	c.page.Timeout(2 * time.Second).WaitStable(300 * time.Millisecond)

	logger.Info(ctx, "🔀 Switched to new tab", zap.String("url", c.GetURL()))
	return true
}

// CloseOtherTabs закрывает все вкладки кроме текущей
func (c *Controller) CloseOtherTabs(ctx context.Context) {
	pages, err := c.browser.Pages()
	if err != nil {
		return
	}

	for _, p := range pages {
		if p != c.page {
			p.Close()
		}
	}
}

// GetTabCount возвращает количество открытых вкладок
func (c *Controller) GetTabCount() int {
	pages, err := c.browser.Pages()
	if err != nil {
		return 1
	}
	return len(pages)
}
