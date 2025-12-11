package action

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
)

// Type вводит текст в поле по селектору
func Type(ctx context.Context, p PageProvider, selector, text string) error {
	logger.Info(ctx, "⌨️ Typing", zap.String("selector", selector), zap.Int("len", len(text)))

	page := p.GetPage()
	elem, err := page.Timeout(10 * time.Second).Element(selector)
	if err != nil {
		return fmt.Errorf("element not found: %s", selector)
	}

	elem.ScrollIntoView()
	if err := elem.Click("left", 1); err != nil {
		logger.Warn(ctx, "⚠️ Click before type failed", zap.Error(err))
	}

	time.Sleep(100 * time.Millisecond)

	if err := elem.SelectAllText(); err == nil {
		elem.MustInput("")
	}

	if err := elem.Input(text); err != nil {
		return fmt.Errorf("input failed: %w", err)
	}

	p.WaitStable(2 * time.Second)
	logger.Info(ctx, "✅ Type completed")
	return nil
}
