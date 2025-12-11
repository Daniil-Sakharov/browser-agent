package browser

import (
	"context"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/browser/action"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/browser/dom"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
	"go.uber.org/zap"
)

// Controller управляет браузером
type Controller struct {
	browser   *rod.Browser
	page      *rod.Page
	timeout   time.Duration
	extractor *dom.Extractor
}

// New создаёт новый контроллер браузера
func New(ctx context.Context, headless bool, userDataDir string, timeoutSec int) (*Controller, error) {
	timeout := time.Duration(timeoutSec) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	l := launcher.New().Headless(headless).Devtools(false)
	if userDataDir != "" {
		l = l.UserDataDir(userDataDir)
	}

	u, err := l.Launch()
	if err != nil {
		return nil, fmt.Errorf("launch: %w", err)
	}

	browser := rod.New().ControlURL(u)
	if err := browser.Connect(); err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	page, err := browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return nil, fmt.Errorf("page: %w", err)
	}

	page.Timeout(timeout).WaitLoad()

	logger.Info(ctx, "✅ Browser initialized", zap.Bool("headless", headless), zap.Duration("timeout", timeout))
	return &Controller{browser: browser, page: page, timeout: timeout, extractor: dom.NewExtractor()}, nil
}

// --- PageProvider interface ---

func (c *Controller) GetPage() *rod.Page        { return c.page }
func (c *Controller) GetTimeout() time.Duration { return c.timeout }
func (c *Controller) WaitStable(timeout time.Duration) {
	c.page.Timeout(timeout).WaitStable(300 * time.Millisecond)
}

// --- Page info ---

func (c *Controller) GetURL() string   { return c.page.MustInfo().URL }
func (c *Controller) GetTitle() string { return c.page.MustInfo().Title }

func (c *Controller) GetHTML(ctx context.Context) (string, error) {
	if c.page == nil {
		return "", fmt.Errorf("page is nil")
	}
	return c.page.Timeout(c.timeout).HTML()
}

func (c *Controller) GetPageContext(ctx context.Context) (*domain.PageContext, error) {
	return c.extractor.ExtractContext(ctx, c.page)
}

// --- Actions (delegate to action package) ---

func (c *Controller) Navigate(ctx context.Context, url string) error {
	return action.Navigate(ctx, c, url)
}

func (c *Controller) Click(ctx context.Context, selector string) error {
	urlBefore := c.GetURL()
	tabsBefore := c.GetTabCount()
	
	if err := action.Click(ctx, c, selector); err != nil {
		return err
	}
	
	time.Sleep(200 * time.Millisecond) // Даём время на навигацию/новую вкладку
	
	// Проверяем новые вкладки
	if c.GetTabCount() > tabsBefore {
		c.SwitchToNewTab(ctx)
		return nil
	}
	
	// Если URL изменился - ждём загрузки
	urlAfter := c.GetURL()
	if urlAfter != urlBefore && urlAfter != "" {
		logger.Info(ctx, "🔄 Navigation detected", zap.String("to", urlAfter))
		c.page.Timeout(c.timeout).WaitLoad()
		c.page.Timeout(2 * time.Second).WaitStable(300 * time.Millisecond)
	}
	
	return nil
}

func (c *Controller) ClickAtPosition(ctx context.Context, x, y int) error {
	return action.ClickAtPosition(ctx, c, x, y)
}

func (c *Controller) Type(ctx context.Context, selector, text string) error {
	return action.Type(ctx, c, selector, text)
}

func (c *Controller) Scroll(ctx context.Context, direction string, amount int) error {
	return action.Scroll(ctx, c, direction, amount)
}

func (c *Controller) PressEnter(ctx context.Context) error {
	return action.PressEnter(ctx, c)
}

// --- DOM (delegate to dom package) ---

func (c *Controller) BuildErrorContext(ctx context.Context, failedSelector string, err error) *domain.ErrorContext {
	return dom.BuildErrorContext(ctx, c, failedSelector, err)
}

func FormatErrorContextMessage(ctx *domain.ErrorContext) string {
	return dom.FormatErrorContextMessage(ctx)
}

// --- Close ---

func (c *Controller) Close(ctx context.Context) error {
	if c.browser != nil {
		logger.Info(ctx, "🚫 Closing browser")
		return c.browser.Close()
	}
	return nil
}
