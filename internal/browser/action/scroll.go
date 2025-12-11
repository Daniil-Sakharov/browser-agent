package action

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
)

// Scroll прокручивает страницу через JavaScript (не зависит от настроек ОС)
func Scroll(ctx context.Context, p PageProvider, direction string, amount int) error {
	logger.Info(ctx, "📜 Scrolling", zap.String("direction", direction), zap.Int("amount", amount))

	if amount == 0 {
		amount = 500
	}

	// JavaScript scroll - работает одинаково на всех ОС
	scrollY := amount
	if direction == "up" {
		scrollY = -amount
	}

	page := p.GetPage()
	js := fmt.Sprintf(`() => {
		window.scrollBy({
			top: %d,
			behavior: 'smooth'
		});
		return {
			scrollY: window.scrollY,
			maxScroll: document.documentElement.scrollHeight - window.innerHeight
		};
	}`, scrollY)

	result, err := page.Eval(js)
	if err != nil {
		return fmt.Errorf("scroll failed: %w", err)
	}

	time.Sleep(400 * time.Millisecond) // Ждём завершения smooth scroll

	logger.Info(ctx, "✅ Scroll completed",
		zap.String("direction", direction),
		zap.Int("scrollY", result.Value.Get("scrollY").Int()),
		zap.Int("maxScroll", result.Value.Get("maxScroll").Int()))
	return nil
}
