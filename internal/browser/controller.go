package browser

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/ysmood/gson"

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
	return action.Click(ctx, c, selector)
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

// --- FindElementsLive - live DOM query ---

func (c *Controller) FindElementsLive(ctx context.Context, query string) (string, error) {
	result, err := c.page.Timeout(5 * time.Second).Eval(findElementsJS)
	if err != nil {
		logger.Error(ctx, "❌ FindElementsLive error", zap.Error(err))
		return "Ошибка поиска.", nil
	}
	elements := result.Value.Arr()
	if len(elements) == 0 {
		return "Элементы не найдены.", nil
	}
	logger.Info(ctx, "🔍 FindElementsLive", zap.Int("found", len(elements)))
	return formatElements(elements), nil
}

func formatElements(elements []gson.JSON) string {
	var out strings.Builder
	out.WriteString("🎯 КЛИКАБЕЛЬНЫЕ ЭЛЕМЕНТЫ:\n\n")
	for i, elem := range elements {
		obj := elem.Map()
		sel := obj["displaySelector"].String()
		if sel == "" {
			continue
		}
		out.WriteString(fmt.Sprintf("%d. %s", i+1, sel))
		if css := obj["cssSelector"].String(); css != "" && css != sel {
			out.WriteString(fmt.Sprintf(" [%s]", css))
		}
		out.WriteString("\n")
	}
	out.WriteString("\n💡 Используй text:ТекстКнопки для клика")
	return out.String()
}

// --- Close ---

func (c *Controller) Close(ctx context.Context) error {
	if c.browser != nil {
		logger.Info(ctx, "🚫 Closing browser")
		return c.browser.Close()
	}
	return nil
}

const findElementsJS = `() => {
	const results = [], seen = new Set(), counts = {};
	function isVisible(el) {
		if (!el) return false;
		const s = window.getComputedStyle(el), r = el.getBoundingClientRect();
		return s.display !== 'none' && s.visibility !== 'hidden' && s.opacity !== '0' && r.width > 0 && r.height > 0 && r.top < window.innerHeight && r.bottom > 0;
	}
	function getSelector(el) {
		const tag = el.tagName.toLowerCase();
		if (el.id && el.id.length < 50) return tag + '#' + el.id;
		const aria = el.getAttribute('aria-label');
		if (aria && aria.length < 100) return tag + '[aria-label="' + aria.replace(/"/g, '\\"') + '"]';
		const tid = el.getAttribute('data-testid');
		if (tid) return tag + '[data-testid="' + tid + '"]';
		if (el.name) return tag + '[name="' + el.name + '"]';
		return tag;
	}
	function getText(el) {
		let t = '';
		for (const n of el.childNodes) if (n.nodeType === Node.TEXT_NODE) t += n.textContent;
		t = t.trim();
		if (t && t.length > 0 && t.length < 50) return t;
		return (el.innerText || '').trim().split('\n')[0].substring(0, 50);
	}
	const all = [], clickable = ['button','a[href]','[role="button"]','[role="menuitem"]','input:not([type="hidden"])','textarea','select','[onclick]','li','label'];
	const modals = ['[role="dialog"]','[aria-modal="true"]','.modal','[class*="modal"]'];
	let hasModal = false;
	for (const ms of modals) { try { for (const m of document.querySelectorAll(ms)) { if (!isVisible(m)) continue; hasModal = true;
		for (const el of m.querySelectorAll(clickable.join(','))) { if (isVisible(el)) { const s = getSelector(el); counts[s] = (counts[s] || 0) + 1; all.push({ el, s }); }}
	}} catch(e){} }
	if (!hasModal) { for (const sel of clickable) { try { for (const el of document.querySelectorAll(sel)) { if (isVisible(el)) { const s = getSelector(el); counts[s] = (counts[s] || 0) + 1; all.push({ el, s }); }}} catch(e){} }}
	for (const item of all) {
		const el = item.el, css = item.s, text = getText(el), r = el.getBoundingClientRect(), key = css + '|' + text;
		if (seen.has(key)) continue; seen.add(key);
		results.push({ displaySelector: text ? 'text:' + text : css, cssSelector: counts[css] === 1 ? css : '', x: Math.round(r.x + r.width / 2), y: Math.round(r.y + r.height / 2) });
		if (results.length >= 60) break;
	}
	return results;
}`
