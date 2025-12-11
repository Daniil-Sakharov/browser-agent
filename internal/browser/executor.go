package browser

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
	"go.uber.org/zap"
)

func (c *Controller) ExecuteAction(ctx context.Context, a domain.Action) (*domain.ActionResult, error) {
	switch a.Type {
	case domain.ActionTypeNavigate:
		return c.exec(ctx, a, func() error { return c.Navigate(ctx, a.URL) }, "Navigated to "+a.URL)
	case domain.ActionTypeClick:
		return c.execWithErr(ctx, a, func() error { return c.Click(ctx, a.Selector) }, "Clicked "+a.Selector)
	case domain.ActionTypeClickAtPosition:
		return c.exec(ctx, a, func() error { return c.ClickAtPosition(ctx, a.X, a.Y) }, fmt.Sprintf("Clicked at (%d,%d)", a.X, a.Y))
	case domain.ActionTypeType:
		return c.execWithErr(ctx, a, func() error { return c.Type(ctx, a.Selector, a.Value) }, "Typed into "+a.Selector)
	case domain.ActionTypeScroll:
		dir := a.Direction; if dir == "" { dir = "down" }
		return c.exec(ctx, a, func() error { return c.Scroll(ctx, dir, 500) }, "Scrolled "+dir)
	case domain.ActionTypeWait:
		return c.execWithErr(ctx, a, func() error { return c.WaitForElement(ctx, a.Selector, 10*time.Second) }, "Element appeared: "+a.Selector)
	case domain.ActionTypePressEnter:
		return c.exec(ctx, a, func() error { return c.PressEnter(ctx) }, "Pressed Enter")
	case domain.ActionTypeCompleteTask:
		return ok(a, "Task completed"), nil
	case domain.ActionTypeTakeScreenshot:
		return c.execScreenshot(ctx, a)
	case domain.ActionTypeQueryDOM:
		return fail(a, "query_dom handled by agent"), nil
	case domain.ActionTypeListTabs:
		return c.execListTabs(ctx, a)
	case domain.ActionTypeSwitchTab:
		return c.exec(ctx, a, func() error { return c.SwitchToTab(ctx, a.TabIndex) }, fmt.Sprintf("Switched to tab %d", a.TabIndex))
	case domain.ActionTypeCloseTab:
		return c.exec(ctx, a, func() error { return c.CloseCurrentTab(ctx) }, "Closed current tab")
	default:
		return nil, fmt.Errorf("unknown action: %s", a.Type)
	}
}

func (c *Controller) exec(ctx context.Context, a domain.Action, fn func() error, msg string) (*domain.ActionResult, error) {
	if err := fn(); err != nil {
		logger.Error(ctx, "❌ Action error", zap.String("action", string(a.Type)), zap.Error(err))
		return fail(a, msg+" failed: "+err.Error()), nil
	}
	return ok(a, msg), nil
}

func (c *Controller) execWithErr(ctx context.Context, a domain.Action, fn func() error, msg string) (*domain.ActionResult, error) {
	if err := fn(); err != nil {
		r := fail(a, msg+" failed")
		r.ErrorContext = c.BuildErrorContext(ctx, a.Selector, err)
		return r, nil
	}
	return ok(a, msg), nil
}

func (c *Controller) execScreenshot(ctx context.Context, a domain.Action) (*domain.ActionResult, error) {
	s, err := c.TakeScreenshot(ctx, a.FullPage, "screenshots")
	if err != nil { return fail(a, "Screenshot failed"), nil }
	r := ok(a, "Screenshot: "+s.Path)
	r.Screenshot, r.ScreenshotB64 = s.Path, s.Base64
	return r, nil
}

func (c *Controller) execListTabs(ctx context.Context, a domain.Action) (*domain.ActionResult, error) {
	result := c.ListTabs(ctx)
	r := ok(a, result)
	r.QueryResult = result
	return r, nil
}

func ok(a domain.Action, msg string) *domain.ActionResult {
	return &domain.ActionResult{Success: true, Action: string(a.Type), Message: msg}
}

func fail(a domain.Action, msg string) *domain.ActionResult {
	return &domain.ActionResult{Success: false, Action: string(a.Type), Message: msg}
}
