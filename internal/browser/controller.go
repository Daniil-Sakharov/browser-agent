package browser

import (
	"context"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/domain"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/logger"
	"go.uber.org/zap"
)

type Controller struct {
	browser *rod.Browser
	page    *rod.Page
	timeout time.Duration
}

type Config interface {
	Headless() bool
	UserDataDir() string
	Timeout() int
}

func New(ctx context.Context, cfg Config) (*Controller, error) {
	timeout := time.Duration(cfg.Timeout()) * time.Second
	if timeout == 0 { timeout = 30 * time.Second }

	l := launcher.New().Headless(cfg.Headless()).Devtools(false)
	if cfg.UserDataDir() != "" { l = l.UserDataDir(cfg.UserDataDir()) }

	u, err := l.Launch()
	if err != nil { return nil, fmt.Errorf("launch: %w", err) }

	browser := rod.New().ControlURL(u)
	if err := browser.Connect(); err != nil { return nil, fmt.Errorf("connect: %w", err) }

	page, err := browser.Page(proto.TargetCreateTarget{})
	if err != nil { return nil, fmt.Errorf("page: %w", err) }

	page.Timeout(timeout).WaitLoad()

	logger.Info(ctx, "✅ Browser initialized", zap.Bool("headless", cfg.Headless()), zap.Duration("timeout", timeout))
	return &Controller{browser: browser, page: page, timeout: timeout}, nil
}

func (c *Controller) Page() *rod.Page { return c.page }

func (c *Controller) GetPageContext(ctx context.Context) (*domain.PageContext, error) {
	return NewExtractor().ExtractContext(ctx, c.page)
}

func (c *Controller) GetHTML(ctx context.Context) (string, error) {
	if c.page == nil { return "", fmt.Errorf("page is nil") }
	return c.page.Timeout(c.timeout).HTML()
}

func (c *Controller) Close(ctx context.Context) error {
	if c.browser != nil {
		logger.Info(ctx, "🚫 Closing browser")
		return c.browser.Close()
	}
	return nil
}
