package action

import (
	"context"

	"github.com/go-rod/rod/lib/input"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
)

// PressEnter нажимает Enter
func PressEnter(ctx context.Context, p PageProvider) error {
	logger.Info(ctx, "⏎ Pressing Enter")
	page := p.GetPage()
	return page.Keyboard.Press(input.Enter)
}

// PressKey нажимает произвольную клавишу
func PressKey(ctx context.Context, p PageProvider, key input.Key) error {
	page := p.GetPage()
	return page.Keyboard.Press(key)
}
