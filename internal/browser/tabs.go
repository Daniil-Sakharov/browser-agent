package browser

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
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

// ListTabs возвращает список всех открытых вкладок
func (c *Controller) ListTabs(ctx context.Context) string {
	pages, err := c.browser.Pages()
	if err != nil {
		return "Ошибка получения списка вкладок"
	}

	var out strings.Builder
	out.WriteString(fmt.Sprintf("📑 Открыто вкладок: %d\n\n", len(pages)))

	currentIdx := -1
	for i, p := range pages {
		if p == c.page {
			currentIdx = i + 1
		}
		info := p.MustInfo()
		title := info.Title
		if len(title) > 50 {
			title = title[:50] + "..."
		}
		marker := "  "
		if p == c.page {
			marker = "→ "
		}
		out.WriteString(fmt.Sprintf("%s%d. %s\n   URL: %s\n", marker, i+1, title, info.URL))
	}

	out.WriteString(fmt.Sprintf("\n💡 Текущая вкладка: %d. Используй switch_tab для переключения", currentIdx))
	logger.Info(ctx, "📑 Listed tabs", zap.Int("count", len(pages)), zap.Int("current", currentIdx))
	return out.String()
}

// SwitchToTab переключается на вкладку по индексу (1-based)
func (c *Controller) SwitchToTab(ctx context.Context, index int) error {
	pages, err := c.browser.Pages()
	if err != nil {
		return fmt.Errorf("failed to get pages: %w", err)
	}

	if index < 1 || index > len(pages) {
		return fmt.Errorf("invalid tab index: %d (available: 1-%d)", index, len(pages))
	}

	c.page = pages[index-1]
	c.page.Timeout(c.timeout).WaitLoad()
	c.page.Timeout(2 * time.Second).WaitStable(300 * time.Millisecond)

	logger.Info(ctx, "🔀 Switched to tab", zap.Int("index", index), zap.String("url", c.GetURL()))
	return nil
}

// CloseCurrentTab закрывает текущую вкладку и переключается на предыдущую
func (c *Controller) CloseCurrentTab(ctx context.Context) error {
	pages, err := c.browser.Pages()
	if err != nil {
		return fmt.Errorf("failed to get pages: %w", err)
	}

	if len(pages) <= 1 {
		return fmt.Errorf("cannot close last tab")
	}

	// Находим индекс текущей вкладки
	currentIdx := 0
	for i, p := range pages {
		if p == c.page {
			currentIdx = i
			break
		}
	}

	// Закрываем текущую
	c.page.Close()

	// Переключаемся на предыдущую или следующую
	newIdx := currentIdx - 1
	if newIdx < 0 {
		newIdx = 0
	}
	c.page = pages[newIdx]
	c.page.Timeout(c.timeout).WaitLoad()

	logger.Info(ctx, "🗑️ Closed tab, switched to", zap.Int("index", newIdx+1), zap.String("url", c.GetURL()))
	return nil
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
